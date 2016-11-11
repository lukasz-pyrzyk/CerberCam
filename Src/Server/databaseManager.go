package main

import "gopkg.in/mgo.v2"

// InsertToDatabase - insert message to database
func InsertToDatabase(msg *Message) {
	session, err := mgo.Dial("MongoDB")
	failOnError(err, "Unable to connect to MongoDB")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("cerberServer").C("messages")
	err = c.Insert(msg)
	failOnError(err, "Unable to insert to database")

	defer session.Close()
}

// ReceiveFromDatabase - Select messages from database
func ReceiveFromDatabase() []Message {
	session, err := mgo.Dial("cerbercam.cloudapp.net")
	failOnError(err, "Unable to connect to MongoDB")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	var msg []Message

	c := session.DB("cerberServer").C("messages")
	err = c.Find(nil).Sort("-_id").Limit(50).All(&msg)

	failOnError(err, "Unable to select from database")

	defer session.Close()

	return msg
}
