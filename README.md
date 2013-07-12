# Gros Canot aka '[Rabaska](https://fr.wikipedia.org/wiki/Rabaska)'

DISCLAIMER: I DON'T ACTUALLY KNOW MY STUFF.

"Gros Canot" means Big Canoe.  It's a play on uOttawa's Rabaska, which is ridiculously named and a real PITA to use.

## Warning

Don't run the `db_viewer.go` file with any `backfill-*` argument.  If you manage to run it anyway, __do not lower the values controlling the delay__ between queries.  If you manage to do that too, __you're solely responsible for whatever stupid thing you do and what might or will happen to you__.

## Creating an API for an University's calendar

I've never done anything REST.  I've never made a REST/JSON API.  I've not read
the REST paper, I've listened to it in my bath while OS X's speech
synthetiser read it on my behalf (thanks for that buddy).

Let's make a REST/JSON API for uOttawa.

## Why am I doing this?

### I'm making an app to test out my key-value store
That's actually the real reason.  I just finished building alpha release of [`dskvs`](https://github.com/aybabtme/dskvs) (star it!) and I need to assert that its API and capabilities are sufficient.

### uOttawa's website is terrible
It's true.

## Tech spec

Let's make a tech spec, because tech spec is a word that rolls in the mouth.
Well it doesn't roll, it more 'clacks' or 'kaplaks' in the mouth.

### Backend

I'm gonna make my backend in Go, because I like Go and thus any other stack for a backend would obviously be a retarded stack.

NoSQL is a fancy way to say, I need a database that has no SQL shit.  Well that
was a pretty obvious one.

I'm fancy and I made my own NoSQL database.  Here it is,
[`dskvs`](https://github.com/aybabtme/dskvs) (yes this is a link). It stands for _Dead Simple Key Value Store_.  I have no
imagination.

### Frontend

I'm gonna make a website facing the API.  It's gonna be built with Angular and
Dart.  Because I'm a Google fanboy.  They make good tech, whatever.

### Data model

I suspect uOttawa's way to organize their educational business is pretty common
in the educational business world.  That is:

* They have programs, or degrees. We'll call them _degrees_ because to me, a _program_ is a binary that happens to be executable.
* Degrees have courses.  Some courses are optional or interchangeable. Some
are mandatory.
* Courses have zero-to-many dependencies to other courses.
* As silly as it sounds, there's no garanty that a degree will be _correct_.
An incorrect degree would, for instance, propose that you take an optional
biology course and restrict you from tacking any of its requirements.

That makes a degree being a forest of DAG.  Wow that was fancy, but what does
that even mean?  Let's draw it:

``` INSERT A PICTURE OF THE WINDOWS ```

I'm going to represent the uOttawa business in JSON.  It's gonna be a lot of
fun. To begin, here's how I see a degree:

```Javascript
{
	"id" : 1209,
	"name" : "Software Engineering",
	"description" : "Learn how to do UML diagrams, code them in Java and lose your soul doing so.",
	"credit" : 132,
	"mandatory_list" : [
		"ITI1100",
		...
		"SEG4110",
	],
	"option" : [
		"GNG4310",
		...,
		"CHG2317",
	],
}
```

And here's a course:
```Javascript
{
	"id" : "ITI1100"
	"topic" : "ITI",
	"code" : "1100",
	"level" : 1000,
	"credit" : 3,
	"name" : "Introduction to Computing",
	"description" : "Learn how to code in Java",
	"dependency" : [],
	"equivalence" : ["ITI1500"],
}
```

On top of that, well need, for various reasons, to hold a representation of what a topic is:

```Javascript
{
	"code": "SEG",
	"descr": "Become a Java dev, draw UML diagrams everyday",
	"courses" : [
		"SEG2105",
		...,
		"SEG4910"
	],
}
```

Ok, that's not really how __I__ see it... but that's how uOttawa has organized it on their website.

## Data gathering

So we need to gather data about _degrees_ and their _courses_, and all the _topics_.

### Degrees
It just so happens that uOttawa lists all their degrees in a convenient
place (that is inconveniently hidden):

http://www.uottawa.ca/academic/info/regist/calendars/programs/

One problem solved.

### Topics

Before I can get a source for the courses data, I first need to get a list of all the topic codes.  A topic code, in my very own newspeak, is the three letters before a course code, like `ITI` in `ITI1100`.

This is the only page that I call that is dynamically generated:

http://www.registrar.uottawa.ca/Default.aspx?tabid=3516

Luckily, I call it only once in the whole scrappe because it contains all the data I need.

### Courses

From the topic codes, I can get a description of every course:

http://www.uottawa.ca/academic/info/regist/calendars/courses/{{TopicCode}}.html

This is nice because all these files are static files and won't make any database call.  I can thus rest assured that I'm not going to DoS uOttawa.ca when I scrappe their thing.

### Putting it all together

So the general idea here is that we will use the list of topics to find all the courses we need to query.  Then one by one, we will request the details for each course.  At the time of writing, there are [5185 available courses](https://web30.uottawa.ca/v3/SITS/timetable/SearchResults.aspx) (this link would be awesome if it wasn't paginated with Javascript links) at uOttawa.

Before doing that, I'll note that the URLs mentionned above are not in the `robots.txt` file of the uottawa.ca website.

```
% curl http://www.uottawa.ca/robots.txt
User-Agent: *
Disallow: /testing/
Disallow: /services/ccs/test/
Disallow: /academic/info/info2/
```

I've setup the scrapper to not query the uOttawa website more than once per 1.5 second, and I only query static endpoints aside for the one call to the topic page.

I've also contacted the uOttawa IT department before scrapping the data and asked them for information about how I could avoid disrupting their operations.  They did not give me any specific, only told me not to query too often as I would get my IP blacklisted by their 'DoS detector'.

### REST endpoints

I need courses, degrees and topics (`dskvs` has no index/ranging mechanism).

```
GET /courses 		# List of all the courses
GET /courses/:id	# Course with the given id
GET /degrees 		# List of all the degrees
GET /degrees/:id 	# A degree with the given id
GET /topics 		# List of all the topics
GET /topics/:id 	# A list of all the courses in that topic
```

Just for fun, I'll prefix all of those endpoints with `v1` and pretend that I have any intention of providing `v2`.

```
/v1/courses
/v1/courses/:id
/v1/degrees
/v1/degrees/:id
/v1/topics
/v1/topics/:id
```

## What's next.

That's it.  Now, all I have to do is:

* Get the data in `dskvs`. __That's done.__
* Build the backend with [Revel](http://robfig.github.io/revel/). __Done as well.__
* Build the frontend with Angular/Dart.
* Become filthy rich!
