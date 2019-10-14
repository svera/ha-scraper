package cmd

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
	"github.com/svera/ha-scraper/providers"
)

var maxDate string
var err error

var rootCmd = &cobra.Command{
	Use:   "<keywords>",
	Short: "Scrapes HA comic images",
	Long:  `Scrapes HA comic images`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		c := colly.NewCollector(
			//colly.Async(true),
			// Cache responses to prevent multiple download of pages
			// even if the collector is restarted
			//colly.CacheDir("./epublibre_cache"),
			colly.AllowURLRevisit(),
		)

		keywords := strings.Join(args, " ")
		defer func() {
			if r := recover(); r != nil {
				log.Printf("run time panic: %v", r)

				// if you just want to log the panic, panic again
				panic(r)
			}
		}()
		start := time.Now()

		log.Printf("Scraping for keywords '%s' started at %02d:%02d:%02d", keywords, start.Hour(), start.Minute(), start.Second())
		files, err := providers.Scrape(c, url.QueryEscape(keywords))
		if err != nil {
			log.Println(err.Error())
		}
		bar := pb.StartNew(len(files))
		log.Printf("Downloading %d files...\n", len(files))
		downloaded := 0
		for url, fileName := range files {
			if err := downloadFile(fileName+".jpg", url); err != nil {
				log.Println(err.Error())
			} else {
				downloaded++
			}
			bar.Increment()
			//log.Printf("Downloaded %s\n", fileName+".jpg")
		}
		bar.Finish()
		end := time.Now()
		log.Printf("Scraping finished at %02d:%02d:%02d, %d images downloaded", end.Hour(), end.Minute(), end.Second(), downloaded)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
