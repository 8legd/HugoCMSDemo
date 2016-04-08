package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/astaxie/beego/session"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"

	"github.com/8legd/hugocms/config"
	hugocms_qor "github.com/8legd/hugocms/qor"
)

var SessionManager *session.Manager

type Auth struct{}

func (Auth) LoginURL(c *admin.Context) string {
	return "/login"
}

func (Auth) LogoutURL(c *admin.Context) string {
	return "/admin/logout"
}

func (Auth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	w := c.Writer
	r := c.Request
	sess, err := SessionManager.SessionStart(w, r)
	if err != nil {
		handleError(err)
	}
	defer sess.SessionRelease(w)

	if r.URL.String() == "/admin/auth" &&
		r.FormValue("inputAccount") != "" &&
		(r.FormValue("inputAccount") == os.Getenv("HUGOCMS_ACC")) &&
		r.FormValue("inputPassword") != "" &&
		(r.FormValue("inputPassword") == os.Getenv("HUGOCMS_PWD")) {
		sess.Set("User", r.FormValue("inputAccount"))
	}
	if userName, ok := sess.Get("User").(string); ok && userName != "" {
		return User{userName}
	}
	return nil
}

type User struct {
	Name string
}

func (u User) DisplayName() string {
	return u.Name
}

func main() {
	var reset bool
	flag.BoolVar(&reset, "reset", false, "reset to an empty database")
	var dbConn string
	flag.StringVar(&dbConn, "db", "", "database connection string")
	flag.Parse()

	var db *gorm.DB
	var err error
	dbName := os.Getenv("HUGOCMS_DBN")
	if dbName == "" {
		handleError(fmt.Errorf("missing env. var. %s", "HUGOCMS_DBN"))
	}
	if dbConn != "" { // to use mysql instead specify a connection string through the db flag
		db, err = gorm.Open("mysql", dbConn+"/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	} else {
		db, err = gorm.Open("sqlite3", dbName+".db")
	}

	if err != nil {
		handleError(err)
	}
	db.LogMode(false)

	if reset {
		// setup empty database
		for _, table := range hugocms_qor.Tables {
			if err := db.DropTableIfExists(table).Error; err != nil {
				handleError(err)
			}
			if err := db.AutoMigrate(table).Error; err != nil {
				handleError(err)
			}
		}

	}

	if err := config.Setup("qor.toml", "hugo.toml", db, Auth{}); err != nil {
		handleError(err)
	}

	// Add session support - used by Auth
	sessionLifetime := 3600 // session lifetime in seconds
	SessionManager, err = session.NewManager("memory", fmt.Sprintf(`{"cookieName":"gosessionid","gclifetime":%d}`, sessionLifetime))
	if err != nil {
		handleError(err)
	}
	go SessionManager.GC()

	// Create Hugo's content directory if it doesnt exist
	// TODO read content dir from config
	if _, err := os.Stat("./content"); os.IsNotExist(err) {
		err = os.MkdirAll("./content", os.ModePerm)
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("public")))

	adm := hugocms_qor.SetupAdmin()

	adm.MountTo("/admin", mux)
	adm.GetRouter().Post("/auth", func(ctx *admin.Context) {
		// we will only hit this on succesful login - redirect to admin dashboard
		w := ctx.Writer
		r := ctx.Request
		http.Redirect(w, r, "/admin", http.StatusFound)
	})
	adm.GetRouter().Get("/logout", func(ctx *admin.Context) {
		w := ctx.Writer
		r := ctx.Request
		sess, err := SessionManager.SessionStart(w, r)
		if err != nil {
			handleError(err)
		}
		defer sess.SessionRelease(w)
		sess.Delete("User")
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	// NOTE: `system` is where QOR admin will upload files e.g. images - we map this to Hugo's static dir along with our other static assets
	// TODO read static dir from config
	// TODO read static assets list from config
	for _, path := range []string{"system", "css", "fonts", "images", "js", "login"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("static")))
	}

	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", config.QOR.Port), mux); err != nil {
		handleError(err)
	}

	fmt.Printf("Listening on: %v\n", config.QOR.Port)
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
	//TODO more graceful exit!
}
