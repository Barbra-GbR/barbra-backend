package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Barbra-GbR/barbra-backend/db"
	"github.com/Barbra-GbR/barbra-backend/server"
	"github.com/Barbra-GbR/barbra-backend/config"
	"github.com/Barbra-GbR/barbra-backend/auth"
	"github.com/Barbra-GbR/barbra-backend/helpers"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	//Initialization
	config.Initialize(*environment)
	helpers.InitializeValidator()
	db.Initialize()
	auth.InitializeJWT()
	auth.InitializeAccountManager()

	//Initialize and start server
	server.Initialize()
}
