package main

import (
	"flag"

	"github.com/8legd/hugocms/server"
)

func main() {
	var reset bool
	flag.BoolVar(&reset, "reset", false, "reset to an empty database")
	//var dbConn string
	//flag.StringVar(&dbConn, "db", "", "database connection string")
	flag.Parse()
	server.ListenAndServe(reset, "")
}
