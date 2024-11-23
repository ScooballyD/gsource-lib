package server

import (
	"html/template"
	"net/http"

	"github.com/ScooballyD/gsource-lib/scrapers"
)

func StartServer(games []scrapers.Game) {
	gamesTmp := template.Must(template.ParseFiles("templates/games.html"))
	discTmp := template.Must(template.ParseFiles("templates/discounts.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gamesTmp.Execute(w, games)
	})
	http.HandleFunc("/discounts", func(w http.ResponseWriter, r *http.Request) {
		discTmp.Execute(w, games)
	})

	http.ListenAndServe(":8080", nil)
}
