package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"

	"github.com/slashformotion/helix-theme-choser/internal"
)

func main() {
	indexHtmlPath := "public/index.html"
	reportPath := "out/report.json"

	btes, err := os.ReadFile(reportPath)
	if err != nil {
		panic(err)
	}
	var r internal.Report
	err = json.Unmarshal(btes, &r)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(indexHtmlPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tmpl := template.Must(template.ParseFiles("website-builder/main.gotmpl"))
	if err := tmpl.Execute(f, r); err != nil {
		log.Fatal(err)
	}
}
