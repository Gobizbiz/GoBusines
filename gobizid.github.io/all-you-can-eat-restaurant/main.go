package main

import (
    "html/template"
    "net/http"
    "log"
)

func main() {
    // File handler untuk file statis (CSS, gambar)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Routes
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/aboutus", aboutHandler)
    http.HandleFunc("/pricelist", priceListHandler)
    http.HandleFunc("/cs", customerServiceHandler)
    http.HandleFunc("/contact", contactHandler)

    // Start server
    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// Home handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index")
}

// About handler
func aboutHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "aboutus")
}

// Pricelist handler
func priceListHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "pricelist")
}

// Customer service handler
func customerServiceHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "cs")
}

// Contact handler
func contactHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "contact")
}

// Function to render templates
func renderTemplate(w http.ResponseWriter, tmpl string) {
    parsedTemplate, _ := template.ParseFiles("templates/" + tmpl + ".html")
    err := parsedTemplate.Execute(w, nil)
    if err != nil {
        log.Println("Error executing template:", err)
    }
}
