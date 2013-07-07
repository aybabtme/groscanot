# Gros Canoe aka '[Rabaska](https://fr.wikipedia.org/wiki/Rabaska)'

## Creating an API for a Unitersity calendar

I've never done anything REST.  I've never made a REST/JSON API.  I've not read
the REST paper, I've listened to it in my bath with while OS X's speech
synthetiser read it on my behalf (thanks for that buddy).

Let's make a REST/JSON API for uOttawa.

## Tech spec

Let's make a tech spec, because tech spec is a word that rolls in the mouth.
Well it doesn't roll, it more 'clacks' or 'kaplak' in the mouth.

### Backend
I'm gonna make my backend in Go, because I like Go and thus any other stack for
a backend would obviously bee a retarded stack.

NoSQL is a fancy way to say, I need a database that has no SQL shit.  Well that
was a pretty obvious one.

I'm fancy and I made my own NoSQL database.  Here it is,
[`dskvs`](https://github.com/aybabtme/dskvs) (yes this is a https link, I'm that
kind of dude!). It stands for _Dead Simple Key Value Store_.  I have no imagination.

### Frontend

I'm gonna make a website facing the API.  It's gonna be built with Angular and Dart.  Because I'm a Google fanboy.  They make good tech, whatever.

### Data model

I suspect uOttawa's way to organize their educational business is pretty common in the educational business world.  That is:

* They have programs, or degrees.
* Programs have courses.  Some courses are optional or interchangeable.
* Courses zero-to-many dependencies to other courses.

That makes a program being a forest of DAG.  Wow that was fancy, but what does that even mean?  Let's draw it:

``` INSERT A PICTURE OF THE WINDOWS ```
