package main

import "gopkg.in/mgo.v2"

// InsertToDatabase - insert message to database
func InsertToDatabase(msg *Message) {
	session, err := mgo.Dial(GlobalConfig.Mongo.Host)
	failOnError(err, "Unable to connect to MongoDB")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(GlobalConfig.Mongo.Database).C(GlobalConfig.Mongo.MessagesTable)
	err = c.Insert(msg)
	failOnError(err, "Unable to insert to database")

	defer session.Close()
}

// ReceiveFromDatabase - Select messages from database
func ReceiveFromDatabase() []Message {
	session, err := mgo.Dial(GlobalConfig.Mongo.Host)
	failOnError(err, "Unable to connect to MongoDB")

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	var msg []Message

	c := session.DB(GlobalConfig.Mongo.Database).C(GlobalConfig.Mongo.MessagesTable)
	err = c.Find(nil).Sort("-_id").Limit(50).All(&msg)

	failOnError(err, "Unable to select from database")

	defer session.Close()

	return msg
}
