package controller

import (
	"html/template"
	"net/http"
)

func DaHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/Wappers.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
