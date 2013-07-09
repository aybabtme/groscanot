package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
	"time"
)

const (
	DEGREE_URL = "http://www.uottawa.ca/academic/info/regist/calendars/programs/"
)

var rDegUrl *regexp.Regexp = regexp.MustCompile("[0-9]+[.]html")

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
	Name      string
	Url       string
	Credit    int
	Mandatory []string
	Extra     []string
}

var (
	sName      = "#pageTitle h1"
	sCredit    = "#pageTitle h1[align=right]"
	sMandatory = ".course td.code span a"
	sExtra     = ".LineFT"
)

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
		log.Printf("%v", degree)
	}

}

func readDegree(degreeRead chan Degree) {
	defer close(degreeRead)

	degreeList := readDegreeUrlList()

	tick := time.NewTicker(time.Millisecond * 1000)
	defer tick.Stop()

	log.Printf("Found %d URLs to degree pages", len(degreeList))
	for _, degreeUrl := range degreeList {

		fmt.Printf("...")
		<-tick.C
		fmt.Printf(" tic! %s\n", degreeUrl)

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
	doc, err := goquery.NewDocument(DEGREE_URL)
	if err != nil {
		log.Printf("Error getting degree list %s: %v", DEGREE_URL[:10], err)
		return nil
	}

	log.Printf("readDegreeUrlList Reading <%s> done in %s\n",
		DEGREE_URL, time.Since(t0))

	var degrees []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if rDegUrl.MatchString(s.Text()) {
			degrees = append(degrees, s.Text())
		}
	})

	return degrees
}

func readDegreePage(degreePage string) (Degree, error) {

	deg := Degree{Url: DEGREE_URL + degreePage}

	t0 := time.Now()

	doc, err := goquery.NewDocument(deg.Url)
	if err != nil {
		log.Printf("Error getting degree doc %s, %v", degreePage, err)
		return deg, err
	}
	log.Printf("readDegreePage Reading <%s> done in %s\n",
		deg.Url, time.Since(t0))

	deg.Name = doc.Find(sName).First().Text()

	deg.Credit, err = strconv.Atoi(doc.Find(sCredit).First().Text())
	if err != nil {
		log.Printf("Couldn't get int our of credit field, %v", err)
	}

	doc.Find(sMandatory).Each(func(i int, s *goquery.Selection) {
		deg.Mandatory = append(deg.Mandatory, s.Text())
	})

	doc.Find(sExtra).Each(func(i int, s *goquery.Selection) {
		deg.Extra = append(deg.Extra, s.Text())
	})

	return deg, nil

}
