package main

import (
	"encoding/base64"
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

/*
	This code is fucking ugly.  Like really fucking ugly.  It's just the fastest piece of shit I needed to get a working data set.
*/

const (
	// In ms
	DEGREE_QUERY_DELAY = 1500
	COURSE_QUERY_DELAY = 1500

	TOPIC_COLL = "topics"
	TOPIC_URL  = "http://www.registrar.uottawa.ca/Default.aspx?tabid=3516"
	// Ignore first one
	S_T_PAIR = "#dnn_ctr6248_HtmlModule_lblContent tbody tr"
	// Index 0
	S_T_VAL = "td"

	DEGREE_COLL   = "degrees"
	DEGREE_URL    = "http://www.uottawa.ca/academic/info/regist/calendars/programs/"
	S_D_NAME      = "#pageTitle h1"
	S_D_CREDIT    = "#pageTitle h1[align=right]"
	S_D_MANDATORY = ".course td.code span a"
	S_D_EXTRA     = ".LineFT"

	COURSE_COLL  = "courses"
	COURSE_URL   = "http://www.uottawa.ca/academic/info/regist/calendars/courses/"
	S_CRS_BOX    = "#crsBox"
	S_CRS_CODE   = ".crsCode"
	S_CRS_TITLE  = ".crsTitle"
	S_CRS_CREDIT = ".crsCredits"
	S_CRS_DESC   = ".crsDesc"
	S_CRS_REQ    = ".crsRestrict"

	CLASS_COLL = "classes"
	CLASS_URL  = "https://web30.uottawa.ca/v3/SITS/timetable/Course.aspx?code="
	S_C_NAME   = "#main-content h2"
)

var (
	rgxDegUrl    *regexp.Regexp = regexp.MustCompile("[0-9]+[.]html")
	rgxCrsCredit *regexp.Regexp = regexp.MustCompile("([0-9]{1}).cr[.]")
	rgxCrsCode   *regexp.Regexp = regexp.MustCompile("[a-zA-Z]{3}[0-9]{4}")
)

type Topic struct {
	Code        string    `json:"code"`
	Description string    `json:"descr"`
	Courses     []string  `json:"courses"`
	LastUpdated time.Time `json:"updated"`
}

type Course struct {
	Id          string    `json:"id"`
	Topic       string    `json:"topic"`
	Code        string    `json:"code"`
	Url         string    `json:"url"`
	Level       int       `json:"level"`
	Credit      int       `json:"credit"`
	Name        string    `json:"name"`
	Description string    `json:"descr"`
	Dependency  []string  `json:"depend"`
	Equivalence []string  `json:"equiv"`
	LastUpdated time.Time `json:"updated"`
}

type Degree struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Credit      int       `json:"credit"`
	Mandatory   []string  `json:"mandat"`
	Extra       []string  `json:"extra"`
	LastUpdated time.Time `json:"updated"`
}

func main() {

	flagCourse := flag.Bool(
		"courses",
		false,
		"print courses in the datastore",
	)

	flagTopic := flag.Bool(
		"topics",
		false,
		"print topics in the datastore",
	)

	flagDegree := flag.Bool(
		"degrees",
		false,
		"print degrees in the datastore",
	)

	valCourse := flag.String(
		"course",
		"",
		"print value of that course",
	)

	valTopic := flag.String(
		"topic",
		"",
		"print value of that topic",
	)

	valDegree := flag.String(
		"degree",
		"",
		"print value of that degree",
	)

	flagDegreeBF := flag.Bool(
		"backfill-degree",
		false,
		"backfill the degrees from the website",
	)

	flagTopicBF := flag.Bool(
		"backfill-topic",
		false,
		"backfill the topics from the website",
	)

	flagCourseBF := flag.Bool(
		"backfill-course",
		false,
		"backfill the courses from the website",
	)
	flag.Parse()

	if !(*flagCourse ||
		*flagTopic ||
		*flagDegree ||
		*flagCourseBF ||
		*flagTopicBF ||
		*flagDegreeBF ||
		*valCourse != "" ||
		*valTopic != "" ||
		*valDegree != "") {

		log.Printf("%v", *flagTopic)

		flag.Usage()
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
			log.Printf("len(c.Id)=%d, c.Id=\"%s\"", len(c.Id), c.Id)
			fmt.Printf("%+v\n", c)
		}
	}

	if *valCourse != "" {
		key := COURSE_COLL + dskvs.CollKeySep + *valCourse
		val, ok, err := store.Get(key)
		if !ok {
			log.Printf("Not found : %v", key)
		} else if err != nil {
			log.Printf("Error %v", err)
		} else {
			log.Printf(string(val))
		}
	}

	if *flagTopic {
		for _, t := range listTopics(store) {
			fmt.Printf("%+v\n", t)
		}
	}

	if *valTopic != "" {
		key := TOPIC_COLL + dskvs.CollKeySep + *valTopic
		val, ok, err := store.Get(key)
		if !ok {
			log.Printf("Not found : %v", key)
		} else if err != nil {
			log.Printf("Error %v", err)
		} else {
			log.Printf(string(val))
		}
	}

	if *flagDegree {
		for _, d := range listDegrees(store) {
			fmt.Printf("%+v\n", d)
		}
	}

	if *valDegree != "" {
		key := DEGREE_COLL + dskvs.CollKeySep + *valDegree
		val, ok, err := store.Get(key)
		if !ok {
			log.Printf("Not found : %v", key)
		} else if err != nil {
			log.Printf("Error %v", err)
		} else {
			log.Printf(string(val))
		}
	}

	if *flagDegreeBF {
		doDegreeBackfill(store)
	}

	if *flagTopicBF {
		doTopicBackfill(store)
	}

	if *flagCourseBF {
		doCourseBackfill(store)
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

func listTopics(s *dskvs.Store) []Topic {
	results, err := s.GetAll(TOPIC_COLL)
	if err != nil {
		log.Printf("Couldn't query back saved topics, %v", err)
		return nil
	}
	var topics []Topic
	for _, b := range results {
		d := Topic{}
		if err := json.Unmarshal(b, &d); err != nil {
			log.Printf("Couldn't unmarshal topics from store, %v", err)
			continue
		}
		topics = append(topics, d)
	}
	return topics
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
		key := DEGREE_COLL + dskvs.CollKeySep + degree.Id

		if err = s.Put(key, b); err != nil {
			log.Printf("Error Putting degree, %v", err)
			return
		}
	}
}

func doTopicBackfill(s *dskvs.Store) {
	topicChan := make(chan Topic)

	go readTopicPage(s, topicChan)

	for topic := range topicChan {
		b, err := json.Marshal(topic)
		if err != nil {
			log.Printf("Couldn't marshal topic, %v", err)
			return
		}
		key := TOPIC_COLL + dskvs.CollKeySep + topic.Code

		if err = s.Put(key, b); err != nil {
			log.Printf("Error Putting topic, %v", err)
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
		key := COURSE_COLL + dskvs.CollKeySep + course.Id

		if err = s.Put(key, b); err != nil {
			log.Printf("Error Putting course, %v", err)
			return
		}
	}
	courses := listCourses(s)
	for _, c := range courses {
		lang, err := strconv.Atoi(string(c.Code[1]))
		if err != nil {
			log.Printf("Couldn't get language digit, %v", err)
			continue
		}
		var equiv string
		if lang < 5 && lang >= 0 {
			equiv = c.Id[:4] + strconv.Itoa(lang+4) + c.Id[5:]
		} else if lang >= 5 && lang < 10 {
			equiv = c.Id[:4] + strconv.Itoa(lang-4) + c.Id[5:]
		} else {
			log.Printf("Invalid lang digit=%d", lang)
			continue
		}
		_, ok, err := s.Get(COURSE_COLL + dskvs.CollKeySep + equiv)
		if err != nil {
			log.Printf("Error getting bilingual equiv, %v", err)
			continue
		}
		if !ok {
			log.Printf("Not bilingual, %v", c.Id)
			continue
		}

		for _, known := range c.Equivalence {
			if known == equiv {
				log.Printf("Already know that one, %v", known)
				continue
			}
		}
		c.Equivalence = append(c.Equivalence, equiv)
		log.Printf("Linking bilingual vs of %s: %s", c.Id, equiv)

		b, err := json.Marshal(c)
		if err != nil {
			log.Printf("Couldn't marshal c, %v", err)
			return
		}
		key := COURSE_COLL + dskvs.CollKeySep + c.Id

		if err = s.Put(key, b); err != nil {
			log.Printf("Error Putting c, %v", err)
			return
		}

	}

	reconcileTopicWithCourses(s)
}

func reconcileTopicWithCourses(s *dskvs.Store) {
	courses := listCourses(s)
	var t Topic
	for _, c := range courses {
		key := TOPIC_COLL + dskvs.CollKeySep + c.Topic
		out, ok, err := s.Get(key)
		if !ok || err != nil {
			log.Printf("Something went wrong, course=%s, ok=%v, err=%v",
				c, ok, err)
			continue
		}
		err = json.Unmarshal(out, &t)
		if err != nil {
			log.Printf("Couldn't unmarshal topic, %v", err)
		}
		for _, known := range t.Courses {
			if known == c.Id {
				log.Printf("Already known by this topic, %s", c.Id)
				continue
			}
		}
		t.Courses = append(t.Courses, c.Id)
		log.Printf("Linking %s to %s", t.Code, c.Id)
		in, err := json.Marshal(t)
		if err != nil {
			log.Printf("Couldn't marshal topic %v, %v", t, err)
			continue
		}
		err = s.Put(key, in)
		if err != nil {
			log.Printf("Couldn't Put, key=%s, err=%v, t=%v", key, err, t)
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

	deg := Degree{Url: DEGREE_URL + degreePage, LastUpdated: time.Now()}

	t0 := time.Now()

	doc, err := goquery.NewDocument(deg.Url)
	if err != nil {
		log.Printf("Error getting degree doc %s, %v", degreePage, err)
		return deg, err
	}
	log.Printf("readDegreePage Reading <%s> done in %s\n",
		deg.Url, time.Since(t0))

	deg.Name = doc.Find(S_D_NAME).First().Text()

	deg.Credit, err = strconv.Atoi(doc.Find(S_D_CREDIT).First().Text())
	if err != nil {
		log.Printf("Couldn't get int our of credit field, %v", err)
	}

	doc.Find(S_D_MANDATORY).Each(func(i int, s *goquery.Selection) {
		deg.Mandatory = append(deg.Mandatory, s.Text())
	})

	doc.Find(S_D_EXTRA).Each(func(i int, s *goquery.Selection) {
		deg.Extra = append(deg.Extra, s.Text())
	})

	deg.Id = base64.StdEncoding.EncodeToString([]byte(deg.Name))

	return deg, nil

}

func readTopicPage(s *dskvs.Store, topicChan chan Topic) {

	t0 := time.Now()

	doc, err := goquery.NewDocument(TOPIC_URL)
	if err != nil {
		log.Printf("Error getting topic doc %s, %v", TOPIC_URL, err)
		return
	}

	log.Printf("readTopicPage Reading <%s> done in %s\n",
		TOPIC_URL, time.Since(t0))

	doc.Find(S_T_PAIR).Each(func(i int, s *goquery.Selection) {
		// Skip the first pair, they're header
		if i == 0 {
			return
		}

		t := Topic{LastUpdated: time.Now()}
		s.Find(S_T_VAL).Each(func(i int, s *goquery.Selection) {
			log.Printf("i=%d Topic = %v", i, s.Text())
			switch i {
			case 0:
				t.Code = s.Children().Text()
			case 1:
				t.Description = s.Text()
			default:
				return
			}
		})
		topicChan <- t

	})
	close(topicChan)
}

func readCourse(s *dskvs.Store, courseRead chan Course) {
	tick := time.NewTicker(time.Millisecond * COURSE_QUERY_DELAY)
	defer tick.Stop()

	topics := listTopics(s)
	for i, topic := range topics {
		fmt.Printf("...")
		<-tick.C
		fmt.Printf(" tick! %d/%d topics\n", i, len(topics))

		courses, err := readCourseFromTopicPage(topic.Code)
		if err != nil {
			log.Printf("Error reading topic code %s, %v", topic.Code, err)
			continue
		}

		for _, c := range courses {
			courseRead <- c
		}
	}
	close(courseRead)
}

func readCourseFromTopicPage(topicCode string) ([]Course, error) {
	target := COURSE_URL + topicCode + ".html"

	t0 := time.Now()

	doc, err := goquery.NewDocument(target)
	if err != nil {
		log.Printf("Error getting topic doc %s, %v", target, err)
		return nil, err
	}

	log.Printf("readCourseFromTopicPage Reading <%s> done in %s\n",
		target, time.Since(t0))

	var courses []Course
	doc.Find(S_CRS_BOX).Each(func(i int, s *goquery.Selection) {
		var id string = s.Find(S_CRS_CODE).Text()
		var topic string = topicCode
		var code string = id[3:]
		var url string = target
		var level int
		var credit int
		var name string = s.Find(S_CRS_TITLE).Text()
		var descr string = s.Find(S_CRS_DESC).Text()
		var depend []string
		var equiv []string

		level, err = strconv.Atoi(string(id[3]))
		if err != nil {
			log.Printf("Error reading course level from id %s, %v", id, err)
			return
		}
		creditStr := rgxCrsCredit.FindString(s.Find(S_CRS_CREDIT).Text())
		if len(creditStr) < 1 {
			log.Printf("No credit for id %d", id)
			return
		} else {
			credit, err = strconv.Atoi(string(creditStr[0]))
			if err != nil {
				log.Printf("Error reading course credit from id %s, %v", id, err)
				return
			}
		}
		depend = rgxCrsCode.FindAllString(s.Find(S_CRS_REQ).Text(), -1)

		c := Course{
			Id:          id,
			Topic:       topic,
			Code:        code,
			Url:         url,
			Level:       level,
			Credit:      credit,
			Name:        name,
			Description: descr,
			Dependency:  depend,
			Equivalence: equiv,
			LastUpdated: time.Now(),
		}

		log.Printf("Read course: %v", c)

		courses = append(courses, c)

	})

	return courses, nil

}
