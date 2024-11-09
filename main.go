package main

import (
    "context"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

var client *mongo.Client

// User struct to map user data
type User struct {
    Name     string `bson:"name" json:"name"`
    Email    string `bson:"email" json:"email"`
    Password string `bson:"password" json:"password"`
}

func main() {
    // Connect to MongoDB
    connectDB()

    // File handler untuk file statis (CSS, gambar)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Routes
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/aboutus", aboutHandler)
    http.HandleFunc("/pricelist", priceListHandler)
    http.HandleFunc("/cs", customerServiceHandler)
    http.HandleFunc("/contact", contactHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/daftar", daftarHandler)

    // Start server
    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// Connect to MongoDB
func connectDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    c, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = c.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = c.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    client = c
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

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        // Ambil input dari form login
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Verifikasi data user dari MongoDB
        collection := client.Database("gobiz").Collection("users")
        var user User
        err := collection.FindOne(context.TODO(), bson.M{"email": email, "password": password}).Decode(&user)
        if err != nil {
            http.Error(w, "Login failed. Invalid email or password.", http.StatusUnauthorized)
            return
        }

        fmt.Fprintf(w, "Welcome, %s!", user.Name)
        return
    }
    renderTemplate(w, "login")
}

// Daftar handler (Register)
func daftarHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        // Ambil input dari form pendaftaran
        name := r.FormValue("name")
        email := r.FormValue("email")
        password := r.FormValue("password")

        // Insert data user baru ke MongoDB
        collection := client.Database("gobiz").Collection("users")
        newUser := User{Name: name, Email: email, Password: password}
        _, err := collection.InsertOne(context.TODO(), newUser)
        if err != nil {
            http.Error(w, "Error saving data.", http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Pendaftaran berhasil, %s!", name)
        return
    }
    renderTemplate(w, "daftar")
}

// Function to render templates
func renderTemplate(w http.ResponseWriter, tmpl string) {
    parsedTemplate, err := template.ParseFiles("templates/" + tmpl + ".html")
    if err != nil {
        log.Println("Error parsing template:", err)
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }
    err = parsedTemplate.Execute(w, nil)
    if err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}
