package db

import (
	"gopkg.in/mgo.v2"
	"time"
	"log"
	"github.com/Barbra-GbR/barbra-backend/config"
)

var db *mgo.Database

//Initializes the database-connection with the data specified in the config
func Initialize() {
	c := config.GetConfig()

	session, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Timeout:  time.Duration(c.GetInt("mongodb.timeout")),
			Addrs:    []string{c.GetString("mongodb.host")},
			Username: c.GetString("mongodb.username"),
			Password: c.GetString("mongodb.password"),
		},
	)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	db = session.DB(c.GetString("database_name"))
}

//Returns the initialized Database. Do not call before calling Initialize!
func GetDB() *mgo.Database {
	return db
}
