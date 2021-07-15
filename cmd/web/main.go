package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abhichvn/bookings/pkg/config"
	"github.com/abhichvn/bookings/pkg/handlers"
	"github.com/abhichvn/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}

	app.TemplateCache = tc
	app.UseCache = false
	// pass template cache to renderer
	render.NewTemplate(&app)

	// pass template cache to handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println("Server is running on port: ", portNumber)
	fmt.Println("URL is accessible at: http://localhost", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
