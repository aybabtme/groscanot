package main

import (
	"github.com/moovweb/gokogiri/html"
	"code.google.com/p/go.net/html"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	DEGREE_URL = "http://www.uottawa.ca/academic/info/regist/calendars/programs/"
)

var degUrlMatch *regexp.Regexp = regexp.MustCompile("[0-9]+[.]html")

type Course struct {
	Id          string
	Topic       string
	Code        string
	Url         string
	Level       int
	Credit      int
	Name        string
	Description string
	Dependency  []Course
	Equivalence []Course
}

type Degree struct {
	Id          string
	Name        string
	Description string
	Url         string
	Credit      int
	Mandatory   []Course
	Option      map[string][]Course
}

func NewDegree()

func main() {
	/*
		store, err := dskvs.Open("./db")
		if err != nil {
			log.Printf("Error opening dskvs: %v", err)
			return
		}
		defer func() {
			err := store.Close()
			if err != nil {
				log.Printf("Error closing dskvs: %v", err)
			}
		}()
	*/

	degreeChan := make(chan Degree)

	go readDegree(degreeChan)

	for degree := range degreeChan {
		log.Printf("Received a degree")
		log.Printf("%v", degree)
	}

}

func readDegree(degreeRead chan Degree) {
	defer close(degreeRead)

	degreeList := readDegreeUrlList()

	log.Printf("Found %d URLs to degree pages", len(degreeList))
	for _, degreeUrl := range degreeList {
		deg, err := readDegreePage(degreeUrl)
		if err != nil {
			log.Printf("Error reading degree page, %v", err)
			return
		}
		degreeRead <- deg
	}
}

func readDegreeUrlList() []string {
	t0 := time.Now()
	response, err := http.Get(DEGREE_URL)
	if err != nil {
		log.Printf("Error getting degree list %s: %v", DEGREE_URL[:10], err)
		return nil
	}
	defer response.Body.Close()
	log.Printf("readDegreeUrlList Reading <%s> done in %s\n",
		DEGREE_URL, time.Since(t0))

	var degreeUrls []string
	root, err := html.Parse(response.Body)
	if err != nil {
		log.Printf("Error parsing response body, %v", err)
		return degreeUrls
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if degUrlMatch.MatchString(attr.Val) {
						degreeUrls = append(degreeUrls, degUrlMatch.FindString(attr.Val))
					}
				}
			}

			if degUrlMatch.MatchString(n.Namespace) {
				degreeUrls = append(degreeUrls, degUrlMatch.FindString(n.Namespace))
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(root)

	return degreeUrls
}



type DegreePage struct {
	"//*[@id="pageTitle"]/div[1]/table/tbody/tr/td[1]/h1"
}

func readDegreePage(degreePage string) (Degree, error) {
	target := DEGREE_URL + degreePage
	t0 := time.Now()
	response, err := http.Get(target)
	if err != nil {
		log.Printf("Error getting degree page %s, %v", degreePage, err)
		return Degree{}, err
	}
	log.Printf("readDegreePage Reading <%s> done in %s\n",
		target, time.Since(t0))

	defer response.Body.Close()



	return deg, err

}
