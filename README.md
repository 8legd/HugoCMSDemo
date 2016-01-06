# HugoCMSDemo

## An experiment with [Hugo](https://github.com/spf13/hugo) and [QOR](https://github.com/qor/qor)

To create an empty sqlite database schema / reset the existing one `$go run main.go -reset`

To run the demo `$go run db/main.go` and go check it out at [http://localhost:8000/admin](http://localhost:8000/admin)

By default Hugo's `staticdir` is configured as `public` (QOR Admin uploads files here e.g. images)

Which means Hugo's `publishdir` needs to be something other than `public`, by default it is configured as `site`

The other Hugo directories are left as per defaults e.g. contentdir is `content`
