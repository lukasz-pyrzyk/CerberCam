package main

import "gopkg.in/mgo.v2"

// InsertToDatabase - insert message to database
func InsertToDatabase(msg *Message) {
	session, err := mgo.Dial("localhost")
	failOnError(err, "Unable to connect to MongoDB")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cerberServer").C("messages")
	err = c.Insert(msg)
	failOnError(err, "Unable to insert to database")

	defer session.Close()
}
