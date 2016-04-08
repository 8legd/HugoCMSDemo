# HugoCMSDemo

## An experiment with [Hugo](https://github.com/spf13/hugo) and [QOR](https://github.com/qor/qor)

To run the demo for the first time `source .env && go run main.go -reset && hugo -w` and wait for the sqlite database to be created - see terminal output...

...then go check it out at [http://localhost:8000/admin](http://localhost:8000/admin)

To run the demo again (without creating the database) `source .env && go run main.go && hugo -w`
