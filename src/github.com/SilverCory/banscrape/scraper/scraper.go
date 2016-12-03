package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

// Scraper is the actual working code.
type Scraper struct {
	BansURL    string
	ResultData []BanEntry
}

// BanEntry containins data about the ban.
type BanEntry struct {
	Placed     time.Time
	Reason     string
	UUID       string
	IssuerUUID string
}

// Create Used to create the actual Scraper
func Create(URL string) *Scraper {

	if strings.HasSuffix(URL, "/") {
		URL = strings.TrimSuffix(URL, "/")
	}

	return &Scraper{
		BansURL:    URL,
		ResultData: make([]BanEntry, 0, 0),
	}
}

// Scrape starts the actual scrape.
func (s *Scraper) Scrape(threads, pages int) {

	if pages < 1 {
		pages = s.getTotalPages()
		if pages < 1 {
			return
		}
	}

	if pages == 1 {
		threads = 1
	}

	if pages <= threads {
		threads = pages + 1
	}

	fmt.Printf("Scraping %d pages using %d threads.\n", pages, threads)

	procs := float64(pages-threads) / float64(threads)

	upRate := int(math.Floor(procs))
	if upRate < 0 {
		upRate = 0
	}

	appendLast := false
	if procs != float64(int64(procs)) {
		appendLast = true
	}

	var completedTasks []int
	completedPages := 0
	min := 0

	for myThread := threads; myThread > 0; myThread-- {

		max := min + upRate
		if myThread == 1 && appendLast {
			max = pages
		}

		fmt.Printf("[%d] Started thread for pages %d to %d (uprate %d)\n", myThread, min, max, upRate)

		go s.getBansFromPages(min, max, myThread, func(compPage, threadId int) {
			if compPage < 1 {
				completedTasks = append(completedTasks, threadId)
				fmt.Printf("Thread %d is completed\n", threadId)
			} else {
				completedPages = completedPages + 1
			}
		})

		min = max + 1

	}

	for {

		completed := 0
		for range completedTasks {
			completed++
		}

		if completed >= threads {
			fmt.Printf("All threads (%d) are completed!\n", threads)
			break
		} else {
			time.Sleep(2 * time.Second)
			fmt.Printf("%.0f%% completed!\n", (float64(completedPages)/float64(pages))*100)
		}

	}

	fmt.Printf("\n\n\nUhhh... Printing %d bans to a json file..\n", len(s.ResultData))

	data, err := json.MarshalIndent(&s.ResultData, "", "\t")
	if err != nil {
		fmt.Println("Unable to marshal json...", err)
		return
	}

	err = ioutil.WriteFile("out.json", data, 0644)
	if err != nil {
		fmt.Println("Unable to write json file...", err)
		return
	}

}

// Appends the ban to the ResultData
func (s *Scraper) addBanEntry(ban BanEntry) {
	s.ResultData = append(s.ResultData, ban)
}
