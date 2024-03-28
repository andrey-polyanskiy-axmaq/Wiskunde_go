package main

import (
	"MVEZ20K/internal/config"
	"MVEZ20K/internal/driver"
	"MVEZ20K/internal/handlers"
	"MVEZ20K/internal/helpers"
	"MVEZ20K/internal/models"
	"MVEZ20K/internal/render"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var ErrorLog *log.Logger

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println("Starting application on port", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err.Error())
	}

	if err := db.SQL.Close(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.User{})

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=postgres user=postgres password=qwerty", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	hostPort := fmt.Sprintf("Connected to database %s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	log.Println(hostPort)
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
