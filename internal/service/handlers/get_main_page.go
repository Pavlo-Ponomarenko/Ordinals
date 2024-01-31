package handlers

import (
	"html/template"
	"net/http"
)

type Params struct {
	Address      string
	Inscriptions []InscriptionView
}

func GetMainPage(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		GetMenuPage(w, r)
		return
	}
	q := InscriptionQ(r)
	entities, _ := q.GetInscriptions(address)
	views := InscriptionEntitiesToViews(entities)
	tmpl, err := template.ParseFiles("templates/Main.html")
	params := Params{
		Address:      address,
		Inscriptions: views,
	}
	err = tmpl.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
