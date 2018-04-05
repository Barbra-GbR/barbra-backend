package main

import (
	"flag"
	"fmt"
	"os"
	"./config"
	"./db"
	"./server"
	"./auth"
)

func main() {
	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	//Initialization
	config.Init(*environment)
	db.Init()
	auth.InitJWT()
	auth.InitUserManager()

	//Initialize and start server
	server.Init()
}