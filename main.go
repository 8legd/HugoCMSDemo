package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/8legd/hugocms/config"
	"github.com/8legd/hugocms/qor"

	"github.com/jinzhu/gorm"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	reset bool
)

func main() {

	flag.BoolVar(&reset, "reset", false, "reset to an empty database")
	flag.Parse()

	if err := config.Parse(); err != nil {
		handleError(err)
	}

	//if db, err := gorm.Open("mysql", "<username>:<password.@tcp(<host>:<port>)/hugocms?charset=utf8&parseTime=True&loc=Local"); err != nil {
	if db, err := gorm.Open("sqlite3", "hugocmsdemo.db"); err != nil {
		panic(err) //TODO be more graceful!
	} else {
		config.QOR.DB = &db
	}
	config.QOR.DB.LogMode(true)

	if reset {
		// setup empty database
		for _, table := range qor.Tables {
			if err := config.QOR.DB.DropTableIfExists(table).Error; err != nil {
				panic(err) //TODO be more graceful!
			}
			if err := config.QOR.DB.AutoMigrate(table).Error; err != nil {
				panic(err) //TODO be more graceful!
			}
		}
	}

	admin := qor.SetupAdmin()

	mux := http.NewServeMux()
	admin.MountTo("/admin", mux)
	// system is where QOR admin will upload files e.g. images
	for _, path := range []string{"css", "fonts", "images", "js", "system"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("public")))
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.QOR.Port), mux); err != nil {
		handleError(err)
	}

	fmt.Printf("Listening on: %v\n", config.QOR.Port)
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	//TODO more graceful exit!
}
