package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/bitphinix/barbra-backend/db"
	"github.com/bitphinix/barbra-backend/server"
	"github.com/bitphinix/barbra-backend/config"
	"github.com/bitphinix/barbra-backend/auth"
	"github.com/bitphinix/barbra-backend/helpers"
)

func main() {
	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	//Initialization
	helpers.InitValidator()
	config.Init(*environment)
	db.Init()
	auth.InitJWT()
	auth.InitUserManager()

	//Initialize and start server
	server.Init()
}
