package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/8legd/hugocms/server"
)

// example server implementation using a single -c flag for configuration
func main() {
	config := ""
	flag.StringVar(&config, "c", "", "config in format account:password@address/database e.g. admin:secret@127.0.0.1:8000/mysql")
	flag.Parse()

	acc := ""
	pwd := ""
	add := ""
	db := ""
	var dbType server.DatabaseType

	parts := strings.Split(config, "@")
	if len(parts) > 1 {
		creds := strings.Split(parts[0], ":")
		if len(creds) > 1 {
			acc = creds[0]
			pwd = creds[1]
		}
		servs := strings.Split(parts[1], "/")
		if len(servs) > 1 {
			add = servs[0]
			db = servs[1]
			switch db {
			case "mysql":
				dbType = server.DB_MySQL
			case "sqlite":
				dbType = server.DB_SQLite
			default:
				fmt.Printf("unsupported database specified through -c flag, expected `mysql` or `sqlite` but recieved `%s`\n", db)
				os.Exit(1)
			}
		}
	}

	if acc == "" || pwd == "" || add == "" {
		fmt.Printf(`invalid or missing -c flag
expected config in format account:password@address/database
but recieved %s:%s@%s/%s
`, acc, pwd, add, db)
		os.Exit(1)
	}

	// check basic requirements...
	if _, err := os.Stat("hugo.toml"); os.IsNotExist(err) {
		fmt.Println("missing Hugo config, please create a `hugo.toml` file for your site")
		os.Exit(1)
	}

	server.ListenAndServe(
		add,
		server.Auth{acc, pwd},
		dbType,
	)
}
