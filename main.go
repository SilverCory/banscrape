package main

import (
	"fmt"
	"runtime"

	"github.com/SilverCory/banscrape/scraper"
)

func main() {

	maxProcs := maxParallelism()
	fmt.Printf("Using %d CPU's\n", maxProcs)

	runtime.GOMAXPROCS(maxProcs)

	beanCraft := scraper.Create("https://www.lemoncloud.org/bans/bans.php")
	beanCraft.Scrape(32, 0)

}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
