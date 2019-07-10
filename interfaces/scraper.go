package interfaces

type Scraper interface {
	Scrape() (int, error)
}
