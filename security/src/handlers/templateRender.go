package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var tmpl *template.Template

func InitTemplates() {
	var err error
	tmpl, err = template.ParseGlob("src/templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

// RenderTemplate rend un template donné avec des données spécifiques
func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
