package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
  "bufio"
  "log"
  "time"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func crawl(url string, foundUrls map[string]bool){
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	// var links []string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				foundUrls[url] = true
			}
		}
	}

}

func main() {
  var seedUrls []string
	foundUrls := make(map[string]bool)

	//read input and store in seedUrls array
  file, err := os.Open("./urls.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      seedUrls = append(seedUrls, scanner.Text())
      // fmt.Println(scanner.Text())
  }
	// fmt.Println(seedUrls)

  //Start timer
  startTime := time.Now()

	// // Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		// fmt.Println(url)
		// fmt.Println(crawl(url))
		crawl(url, foundUrls)
		// fmt.Println(linksReturn)
	}
  //
	// Subscribe to both channels
  elapseTime := time.Now().Sub(startTime)

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

  fmt.Println("Total time: ", elapseTime)

}
