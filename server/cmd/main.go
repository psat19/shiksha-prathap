package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	models "github.com/psat/shiksha-prathap/pkg/models"
	"github.com/psat/shiksha-prathap/pkg/models/pg"
)

type User interface {
	Insert(email, password string) (int, error)
	Authenticate(email, password string) (int, error)
	Get(id int) (*models.User, error)
	Update(name, phone string, age int32, id int) error
}

type application struct {
	InfoLog       *log.Logger
	ErrLog        *log.Logger
	TemplateCache map[string]*template.Template
	Session       *sessions.CookieStore
	User          User
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	serverPort := fmt.Sprintf(":%s", os.Getenv("SRV_PORT"))
	secret := os.Getenv("COOKIE_SECRET")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, _ := newTemplateCache("./ui/html/")

	store := sessions.NewCookieStore([]byte(secret))
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode

	app := &application{
		InfoLog:       infoLog,
		ErrLog:        errLog,
		TemplateCache: templateCache,
		Session:       store,
	}

	db, err := app.openDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app.User = &pg.UserModel{DB: db}

	srv := http.Server{
		Addr:         serverPort,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Starting server on port %s\n", serverPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (app application) openDB(connStr string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 40
	retryInterval := 5 * time.Second

	time.Sleep(retryInterval)

	app.InfoLog.Println("Attempting to connect to PostgreSQL...")

	// Robust connection loop to wait for the database service to be ready
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				app.InfoLog.Println("Successfully connected to the database!")
				break
			}
		}

		app.InfoLog.Printf("DB connection failed (Attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to connect to the database after %d attempts", maxRetries)
		}
	}

	return db, nil
}
