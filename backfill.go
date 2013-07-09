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
	// In ms
	DEGREE_QUERY_DELAY = 1000
	COURSE_QUERY_DELAY = 1000

	DEGREE_COLL = "degrees"
	DEGREE_URL  = "http://www.uottawa.ca/academic/info/regist/calendars/programs/"
	S_NAME      = "#pageTitle h1"
	S_CREDIT    = "#pageTitle h1[align=right]"
	S_MANDATORY = ".course td.code span a"
	S_EXTRA     = ".LineFT"

	COURSE_COLL = "courses"
	COURSE_URL  = "https://web30.uottawa.ca/v3/SITS/timetable/Course.aspx?code="
)

var (
	flagCourse   *bool
	flagDegree   *bool
	flagDegreeBF *bool
	flagCourseBF *bool
	rgxDegUrl    *regexp.Regexp = regexp.MustCompile("[0-9]+[.]html")
)

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

func main() {

	flagCourse = flag.Bool("courses", false, "print courses in the datastore")
	flagDegree = flag.Bool("degrees", false, "print degrees in the datastore")
	flagDegreeBF = flag.Bool("backfill-degree", false, "backfill the degrees from the website")
	flagCourseBF = flag.Bool("backfill-course", false, "backfill the courses from the website")
	flag.Parse()

	if !*flagCourse && !*flagDegree && !*flagDegreeBF {
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
		for _, c := range listCourses(store) {
			fmt.Printf("%+v\n", c)
		}
	}

	if *flagDegree {
		listDegrees(store)
		for _, d := range listDegrees(store) {
			fmt.Printf("%+v\n", d)
		}

	}

	if *flagDegreeBF {
		doDegreeBackfill(store)
	}

}

func listCourses(s *dskvs.Store) []Course {
	results, err := s.GetAll(COURSE_COLL)
	if err != nil {
		log.Printf("Couldn't query back saved courses, %v", err)
		return nil
	}
	var courses []Course
	for _, b := range results {
		c := Course{}
		if err := json.Unmarshal(b, &c); err != nil {
			log.Printf("Couldn't unmarshal courses from store, %v", err)
			continue
		}
		courses = append(courses, c)
	}
	return courses
}

func listDegrees(s *dskvs.Store) []Degree {
	results, err := s.GetAll(DEGREE_COLL)
	if err != nil {
		log.Printf("Couldn't query back saved degrees, %v", err)
		return nil
	}
	var degrees []Degree
	for _, b := range results {
		d := Degree{}
		if err := json.Unmarshal(b, &d); err != nil {
			log.Printf("Couldn't unmarshal degrees from store, %v", err)
			continue
		}
		degrees = append(degrees, d)
	}
	return degrees
}

func doDegreeBackfill(s *dskvs.Store) {
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

func doCourseBackfill(s *dskvs.Store) {
	courseChan := make(chan Course)

	go readCourse(s, courseChan)

	for course := range courseChan {
		b, err := json.Marshal(course)
		if err != nil {
			log.Printf("Couldn't marshal course, %v", err)
			return
		}
		key := COURSE_COLL + dskvs.CollKeySep + course.Name

		err = s.Put(key, b)
		if err != nil {
			log.Printf("Error Putting course, %v", err)
			return
		}
	}
}

func readDegree(degreeRead chan Degree) {
	defer close(degreeRead)

	degreeList := readDegreeUrlList()

	tick := time.NewTicker(time.Millisecond * DEGREE_QUERY_DELAY)
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
		if rgxDegUrl.MatchString(s.Text()) {
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

func readCourse(s *dskvs.Store, courseRead chan Course) {
	tick := time.NewTicker(time.Millisecond * COURSE_QUERY_DELAY)
	defer tick.Stop()

	degs := listDegrees(s)
	for i, d := range degs {
		for j, code := range d.Mandatory {
			fmt.Printf("...")
			<-tick.C
			fmt.Printf(" tick! %d/%d degree, %d/%d course in this degree\n",
				i, len(degs), j, len(d.Mandatory))

			if c, ok := canGetFromStore(s, code); ok {
				courseRead <- c
				continue
			}

			c, err := readCoursePage(code)
			if err != nil {
				log.Printf("Error reading course code %s, %v", code, err)
				continue
			}
			courseRead <- c
		}
	}
}

func canGetFromStore(s *dskvs.Store, code string) (Course, bool) {
	c := Course{}
	b, _ := s.Get(code)
	if b != nil {

		err := json.Unmarshal(b, &c)
		if err != nil {
			log.Printf("Couldn't unmarshal saved course, will read it again, %v", err)
			return c, false
		}

		return c, true

	}

	return c, false
}

func readCoursePage(courseCode string) (Course, error) {
	// Stuff we already know about

	lvl, err := strconv.Atoi(string(courseCode[3]))
	if err != nil {
		log.Printf("Couldn't get level from course code, course code must be invalid, %v", err)
		return Course{}, err
	}

	c := Course{
		Id:    courseCode,
		Url:   COURSE_URL + courseCode,
		Topic: courseCode[:3],
		Code:  courseCode[3:],
		Level: lvl * 1000,
	}

	// Stuff we need to find out
	var credit int
	var name string
	var description string
	var dependency []Course
	var equivalence []Course

	c.Credit = credit
	c.Name = name
	c.Description = description
	c.Dependency = dependency
	c.Equivalence = equivalence

	return c, err
}
