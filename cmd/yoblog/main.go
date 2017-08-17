package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/tthanh/yoblog/postgres"
	"github.com/tthanh/yoblog/service"
)

var schema = `
CREATE TABLE IF NOT EXISTS account (
	id CHAR(36) PRIMARY KEY,
	email VARCHAR(50) UNIQUE,
	name VARCHAR(50) NOT NULL,
	created_at INTEGER,
	updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS post (
	id CHAR(36) PRIMARY KEY,
	owner_id CHAR(36) REFERENCES account (id),
	title VARCHAR(254),
	content TEXT,
	created_at INTEGER,
	updated_at INTEGER
)
`

func main() {
	db, err := sqlx.Connect("postgres", "host=127.0.0.1 user=postgres password=123456 dbname=postgres sslmode=disable ")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = db.Close()
	}()

	db.MustExec(schema)

	accountStore := postgres.NewAccountStore(db)

	oauth2ClientID := os.Getenv("OAUTH2_CLIENT_ID")
	oauth2ClientSecret := os.Getenv("OAUTH2_CLIENT_SECRET")
	oauth2RedirectURL := os.Getenv("OAUTH2_REDIRECT_URL")
	oauth2Scopes := strings.Split(os.Getenv("OAUTH2_SCOPE"), ",")
	oauth2State := os.Getenv("OAUTH2_STATE")

	oauth2Config := &oauth2.Config{
		ClientID:     oauth2ClientID,
		ClientSecret: oauth2ClientSecret,
		RedirectURL:  oauth2RedirectURL,
		Scopes:       oauth2Scopes,
		Endpoint:     facebook.Endpoint,
	}

	srv, err := service.New(
		service.SetAccountStore(accountStore),
		service.SetCookieStore([]byte("secret")),
		service.SetCookieName("yoblog"),
		service.SetOAuth2Config(oauth2Config),
		service.SetOAuth2State(oauth2State),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", srv.IndexHandler).Methods("GET")
	r.HandleFunc("/login", srv.LoginHandler).Methods("GET")
	r.HandleFunc("/callback", srv.CallbackHandler).Methods("GET")
	r.HandleFunc("/logout", srv.LogoutHandler).Methods("GET")

	httpSrv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(httpSrv.ListenAndServe())
}
