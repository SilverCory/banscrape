package main

import (
	"fmt"
	"runtime"

	"github.com/SilverCory/banscrape/scraper"
)

func main() {

	// Max out our procs because we really are hard working...?
	maxProcs := maxParallelism()
	fmt.Printf("Using %d CPU's\n", maxProcs)
	runtime.GOMAXPROCS(maxProcs)

	// Create and scrape.
	beanCraft := scraper.Create("https://www.lemoncloud.org/bans/")
	beanCraft.Scrape(32, 0)

}

// Stolen code from somewhere..
func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
