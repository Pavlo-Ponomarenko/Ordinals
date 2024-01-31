package handlers

import (
	"html/template"
	"net/http"
)

type MenuParams struct{}

func GetMenuPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Menu.html")
	err = tmpl.Execute(w, MenuParams{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
