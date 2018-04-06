package db

import (
	"gopkg.in/mgo.v2"
	"time"
	"log"
	"github.com/bitphinix/babra_backend/config"
)

var db *mgo.Database

func Init() {
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

func GetDB() *mgo.Database {
	return db
}
