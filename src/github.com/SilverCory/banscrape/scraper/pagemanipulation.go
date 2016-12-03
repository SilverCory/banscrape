package scraper

import (
	"errors"
	"fmt"
	"strconv"

	"time"

	"net/url"

	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Get bans form the page.. bad code warning.
func (s *Scraper) getBansFromPages(min, max, threadID int, done func(int, int)) {

	if min > max || max <= 0 {
		fmt.Printf("[%d] Min max mismatch! Min: %d, Max: %d\n", threadID, min, max)
		done(-1, threadID)
		return
	}

	if min <= 0 {
		min = 1
	}

	for page := min; page <= max; page++ {

		doc, err := goquery.NewDocument(s.BansURL + "/bans.php?page=" + strconv.Itoa(page))
		if err != nil {
			fmt.Printf("[%d] An error occurred loading page: %d, %q", threadID, page, err)
			done(page, threadID)
			continue
		}

		doc.Find("body > div > div:nth-child(2) > div > table > tbody").Children().Each(s.processPage)
		done(page, threadID)

	}

	done(-1, threadID)

}

// Process the page.
func (s *Scraper) processPage(i int, selection *goquery.Selection) {

	href, exists := selection.Find("td:nth-child(1) > a").Attr("href")
	if !exists {
		return
	}

	s.getUserData(href)

}

// Get user data from the href.
func (s *Scraper) getUserData(href string) {

	if !strings.HasPrefix(href, "/") {
		href = "/" + href
	}

	doc, err := goquery.NewDocument(s.BansURL + href)
	if err != nil {
		doc.Text()
		fmt.Println("An error occurred user page: "+href+". ", err)
		return
	}

	issuer, err := findUUID(doc.Find("body > div > div:nth-child(2) > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > a"))
	if err != nil || issuer == "" {
		issuer = "CONSOLE"
	}

	banned, err := findUUID(doc.Find("body > div > div:nth-child(2) > div > table > tbody > tr:nth-child(1) > td:nth-child(2) > a"))
	if err != nil {
		fmt.Printf("Unable to find UUID for %q", s.BansURL+href)
		return
	}

	reason := doc.Find("body > div > div:nth-child(2) > div > table > tbody > tr:nth-child(3) > td:nth-child(2)").Text()
	date := doc.Find("body > div > div:nth-child(2) > div > table > tbody > tr:nth-child(4) > td:nth-child(2)").Text()

	banTime, err := time.Parse("January 2, 2006, 04:05 PM", date)
	if err != nil {
		banTime = time.Now()
	}

	banentry := BanEntry{
		IssuerUUID: issuer,
		UUID:       banned,
		Reason:     reason,
		Placed:     banTime,
	}

	s.addBanEntry(banentry)

}

// Get the total number of pages.
func (s *Scraper) getTotalPages() int {

	doc, err := goquery.NewDocument(s.BansURL + "/bans.php")
	if err != nil {
		doc.Text()
		fmt.Println("An error occurred loading the home page.", err)
		return -1
	}

	text := doc.Find("body > div > div:nth-child(2) > div > div:nth-child(6) > div").Text()

	pageNumber := strings.Split(text, "/")
	ret, err := strconv.Atoi(pageNumber[len(pageNumber)-1])
	if err != nil {
		panic(err)
	}

	return ret

}

// Finds the UUID from the anchor tag.
func findUUID(UUIDContainer *goquery.Selection) (string, error) {

	href, exists := UUIDContainer.Attr("href")
	if !exists {
		return "", errors.New("Href doesn't exist!")
	}

	val, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return val.Query().Get("uuid"), nil

}
