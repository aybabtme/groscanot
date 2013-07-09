package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aybabtme/dskvs"
	"log"
	"regexp"
	"strconv"
	"time"
)

const (
	DEGREE_COLL = "degrees"

	DEGREE_URL  = "http://www.uottawa.ca/academic/info/regist/calendars/programs/"
	S_NAME      = "#pageTitle h1"
	S_CREDIT    = "#pageTitle h1[align=right]"
	S_MANDATORY = ".course td.code span a"
	S_EXTRA     = ".LineFT"

	COURSE_COLL = "courses"
)

var rDegUrl *regexp.Regexp = regexp.MustCompile("[0-9]+[.]html")

type Course struct {
	Id          string   `json:"id"`
	Topic       string   `json:"topic"`
	Code        string   `json:"code"`
	Url         string   `json:"url"`
	Level       int      `json:"level"`
	Credit      int      `json:"credit"`
	Name        string   `json:"name"`
	Description string   `json:"descr"`
	Dependency  []Course `json:"depend"`
	Equivalence []Course `json:"equiv"`
}

type Degree struct {
	Name      string   `json:"name"`
	Url       string   `json:"url"`
	Credit    int      `json:"credit"`
	Mandatory []string `json:"mandat"`
	Extra     []string `json:"extra"`
}

var (
	flagCourse   = flag.Bool("courses", false, "print courses in the datastore")
	flagDegree   = flag.Bool("degrees", false, "print degrees in the datastore")
	flagBackfill = flag.Bool("backfill", false, "backfill the data from the website")
)

func main() {
	flag.Parse()

	if !*flagCourse && !*flagDegree && !*flagBackfill {
		flag.PrintDefaults()
		return
	}

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

	if *flagCourse {
		listCourses(store)
	}

	if *flagDegree {
		listDegrees(store)
	}

	if *flagBackfill {
		doBackfill(store)
	}

}

func listCourses(s *dskvs.Store) {
	results, err := s.GetAll(COURSE_COLL)
	if err != nil {
		log.Printf("Couldn't query back saved degrees, %v", err)
		return
	}
	for _, b := range results {
		d := Degree{}
		if err := json.Unmarshal(b, &d); err != nil {
			log.Printf("Couldn't unmarshal degrees from store, %v", err)
			return
		}
		fmt.Printf("%+v\n", d)
	}
}

func listDegrees(s *dskvs.Store) {
	results, err := s.GetAll(DEGREE_COLL)
	if err != nil {
		log.Printf("Couldn't query back saved degrees, %v", err)
		return
	}
	for _, b := range results {
		d := Degree{}
		if err := json.Unmarshal(b, &d); err != nil {
			log.Printf("Couldn't unmarshal degrees from store, %v", err)
			return
		}
		fmt.Printf("%+v\n", d)
	}
}

func doBackfill(s *dskvs.Store) {
	degreeChan := make(chan Degree)

	go readDegree(degreeChan)

	for degree := range degreeChan {
		b, err := json.Marshal(degree)
		if err != nil {
			log.Printf("Couldn't marshal degree, %v", err)
			return
		}
		key := DEGREE_COLL + dskvs.CollKeySep + degree.Name

		err = s.Put(key, b)
		if err != nil {
			log.Printf("Error Putting degree, %v", err)
			return
		}
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

	deg.Name = doc.Find(S_NAME).First().Text()

	deg.Credit, err = strconv.Atoi(doc.Find(S_CREDIT).First().Text())
	if err != nil {
		log.Printf("Couldn't get int our of credit field, %v", err)
	}

	doc.Find(S_MANDATORY).Each(func(i int, s *goquery.Selection) {
		deg.Mandatory = append(deg.Mandatory, s.Text())
	})

	doc.Find(S_EXTRA).Each(func(i int, s *goquery.Selection) {
		deg.Extra = append(deg.Extra, s.Text())
	})

	return deg, nil

}
