# Gros Canoe aka '[Rabaska](https://fr.wikipedia.org/wiki/Rabaska)'

## Creating an API for an University's calendar

I've never done anything REST.  I've never made a REST/JSON API.  I've not read
the REST paper, I've listened to it in my bath while OS X's speech
synthetiser read it on my behalf (thanks for that buddy).

Let's make a REST/JSON API for uOttawa.

## Tech spec

Let's make a tech spec, because tech spec is a word that rolls in the mouth.
Well it doesn't roll, it more 'clacks' or 'kaplaks' in the mouth.

### Backend
I'm gonna make my backend in Go, because I like Go and thus any other stack for a backend would obviously be a retarded stack.

NoSQL is a fancy way to say, I need a database that has no SQL shit.  Well that
was a pretty obvious one.

I'm fancy and I made my own NoSQL database.  Here it is,
[`dskvs`](https://github.com/aybabtme/dskvs) (yes this is a https link, I'm that
kind of dude!). It stands for _Dead Simple Key Value Store_.  I have no
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
		{
			"count": 6,
			"selection" : ["GNG4310", "CHG2317"],
		},
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

Ok, that's not really how __I__ see it... but that's how uOttawa has organized it on their website.

## Data gathering

So we need to gather data about _degrees_ and their _courses_.

### Degrees
It just so happens that uOttawa lists all their degrees in a convenient
place (that is inconveniently hidden):

http://www.uottawa.ca/academic/info/regist/calendars/programs/

One problem solved.

### Courses

From every degree, it's possible to obtain a list of courses code.

A little digging on the uOttawa website reveals that the following URL can be
used to gather information specific to a single course:

https://web30.uottawa.ca/v3/SITS/timetable/Course.aspx?code={{CourseCode}}

We can also use this link to get more details about the instances of the
courses, like when are they given during term X, when are the classes/labs/etc.

### Put it all together

So the general idea here is that we will use the list of degrees to find all the courses we need to query.  Then one by one, we will request the details for each course.  At the time of writing, there are [5185 available courses](https://web30.uottawa.ca/v3/SITS/timetable/SearchResults.aspx) (this link would be awesome if it wasn't paginated with Javascript links) at uOttawa.

Before doing that, I'll note that the URLs mentionned above are not in the `robots.txt` file of the uottawa.ca website.

```
% curl http://www.uottawa.ca/robots.txt
User-Agent: *
Disallow: /testing/
Disallow: /services/ccs/test/
Disallow: /academic/info/info2/
```

I've also attempted to contact the responsible of their IT department, seeking info about their capabilities so that I don't inadvertebly DoS their website. I'll wait for an answer for a week before starting to crawl.

Other than that, I'm not going to query their stuff with 100 concurrent connections.  I think this is a civilized way to do this.

### REST endpoints

I need courses, degrees, maybe topics (`dskvs` has no index/ranging mechanism).

```
GET /course 		# List of all the courses
GET /course/:id		# Course with the given id
GET /degree 		# List of all the degrees
GET /degree/:id 	# A degree with the given id
GET /topic 			# List of all the topics
GET /topic/:id 		# A list of all the courses in that topic
```

Just for fun, I'll prefix all of those endpoints with `v1` and pretend that I have any intention of providing `v2`.

```
/v1/course
/v1/course/:id
/v1/degree
/v1/degree/:id
/v1/topic
/v1/topic/:id
```

## What's next.

That's it.  Now, all I have to do is:

* Get the data in `dskvs`.
* Build the backend with [Revel](http://robfig.github.io/revel/)
* Build the frontend with Angular/Dart.
* Become filthy rich!

## Why am I doing this?

### uOttawa's website sucks
It's true.

### I'm making an app to test out my key-value store
That's actually the real reason.  I just finished building alpha release of `dskvs` and I need to assert that its API and capabilities are sufficient.
