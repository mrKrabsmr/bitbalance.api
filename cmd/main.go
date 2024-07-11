package main

import (
	_ "fl/my-portfolio/docs"
	server "fl/my-portfolio/internal"

	"flag"
)

var options = map[uint8]string{
	1: "migrations-up",
	2: "migrations-down",
}

var opt = flag.String("o", "", "options")

// @title My-Portfolio
// @version 1.0
// @description API for get cryptocurrencies
// @host 127.0.0.1:8000
// @contact.name My-Portfolio support
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	flag.Parse()
	apiServer := server.NewAPIServer()

	if *opt != "" {
		switch *opt {
		case options[1]:
			apiServer.StartMigrations()
		case options[2]:
			apiServer.DownMigrations()
		default:
			panic("unknown option")
		}
		return
	}

	apiServer.Run()
}
