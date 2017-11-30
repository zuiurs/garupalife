package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

const (
	BaseURL            = "https://bang-dream.bushimo.jp"
	ContentsPerPage    = 10
	BackNumberQueryKey = "4koma"
)

var (
	dir       string
	offset    int
	sleepTime int
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "scrape four-panel comic",
	RunE:  scrape,
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringVarP(&dir, "directory", "d", ".", "Destination directory")
	scrapeCmd.Flags().IntVar(&offset, "offset", 1, "Backnumber offset (1 counting)")
	scrapeCmd.Flags().IntVarP(&sleepTime, "time", "t", 500, "Scraping interval in ms")
}

func scrape(cmd *cobra.Command, args []string) error {
	// ---- parameter check ----
	if offset < 1 {
		return fmt.Errorf("Invalid offset: %d", offset)
	}

	dir = strings.TrimSuffix(dir, "/")
	_, err := os.Stat(dir)

	if err != nil {
		return fmt.Errorf("Directory is not found: %s", dir)
	}
	// -------------------------

	imageURLs, err := scrapeAllPages(offset)
	if err != nil {
		return err
	}

	err = fetchAllImages(dir, offset, imageURLs)
	if err != nil {
		return err
	}

	fmt.Println("Complete saving all images")
	return nil
}

func scrapeAllPages(offset int) ([]string, error) {
	imageURLs := make([]string, 0, 100)

	req, err := http.NewRequest("GET", BaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36")

	// start (e.g. 1, 11, 21, 31, ...)
	// Skip disused pages
	start := ((offset-1)/ContentsPerPage)*10 + 1
	for i := start; ; i += ContentsPerPage {
		fmt.Printf("Extract URLs: %d - %d\n", i, i+(ContentsPerPage-1))
		length := len(imageURLs)

		values := url.Values{}
		values.Add(BackNumberQueryKey, fmt.Sprintf("%d_%d", i, i+(ContentsPerPage-1)))
		req.URL.RawQuery = values.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			return nil, err
		}

		doc.Find("div[class=contents_inner_item] > p[class=btn_detail] > a").Each(func(_ int, s *goquery.Selection) {
			if url, ok := s.Attr("href"); ok {
				s, err := findImageURL(url)
				if err != nil {
					return
				}
				imageURLs = append(imageURLs, s)
			}
		})

		// empty page check
		if length == len(imageURLs) {
			break
		}
	}
	return imageURLs, nil
}

func findImageURL(url string) (string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	s, ok := doc.Find("div[class=main_4koma] > p > img").Attr("src")
	if !ok {
		return "", fmt.Errorf("Failed to fetch: %s", url)
	}

	// Filename: c2989ad24ec220478211485d4947f46f-334x1024.png
	//        -> c2989ad24ec220478211485d4947f46f.png
	// Remove hyphend resolution
	hyphenPos := strings.LastIndex(s, "-")
	dotPos := strings.LastIndex(s, ".")
	return fmt.Sprintf("%s%s", s[:hyphenPos], s[dotPos:]), nil
}

func fetchAllImages(dir string, offset int, urls []string) error {
	var c int
	for i, url := range urls {
		if i < (offset-1)%ContentsPerPage {
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		f, err := os.Create(fmt.Sprintf("%s/%d_motto_garupa_life.jpg", dir, offset+c))
		if err != nil {
			return err
		}

		fmt.Printf("Saving Image: %s\n", f.Name())
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
		c++
	}

	return nil
}
