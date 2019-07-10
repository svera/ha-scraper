package providers

import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/gosimple/slug"
)

func Scrape(c *colly.Collector, keywords string) (map[string]string, error) {
	total := 0

	files := map[string]string{}

	c.OnHTML("img.thumbnail.preview", func(e *colly.HTMLElement) {
		total++
		if e.Attr("data-src") == "" {
			return
		}
		decodedValue, err := url.QueryUnescape(e.Attr("data-src"))
		re := regexp.MustCompile(`(,sizedata\[[0-9]+x[0-9]+\])`)
		s := re.ReplaceAllString(decodedValue, "")
		if err != nil {
			log.Println(err)
		}
		text := slug.MakeLang(e.Attr("alt"), "en")
		files[s] = text
		fmt.Printf("%s: %s\n", s, text)
	})

	c.OnHTML("a.icon-right-triangle", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, ". Error", err.Error(), "Retrying...")
	})

	c.Visit(fmt.Sprintf("https://comics.ha.com/c/search-results.zx?N=790+231+52&Ntt=%s&Ntk=SI_Titles&&ic10=AllCategoriesResults-ViewAll-050517", keywords))
	c.Wait()
	return files, nil
}
