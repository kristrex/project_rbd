package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rbd_project/internal/db"
	"rbd_project/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.DB.Close()

	err = db.DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Защищенные эндпоинты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/persons", handlers.GetPersons).Methods("GET")
	api.HandleFunc("/persons", handlers.CreatePerson).Methods("POST")
	api.HandleFunc("/persons", handlers.UpdatePerson).Methods("PUT")
	api.HandleFunc("/persons", handlers.DeletePerson).Methods("DELETE")
	api.HandleFunc("/persons/search", handlers.SearchPersons).Methods("GET")
	api.HandleFunc("/persons/average-age", handlers.GetAverageAgeByCity).Methods("GET")
	api.HandleFunc("/persons/delete-older-than", handlers.DeletePersonsOlderThan).Methods("DELETE")
	api.HandleFunc("/persons/update-with-trigger", handlers.UpdatePersonWithTrigger).Methods("PUT")
	api.HandleFunc("/pizzerias", handlers.GetPizzerias).Methods("GET")
	api.HandleFunc("/visits", handlers.GetVisits).Methods("GET")
	api.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
	api.HandleFunc("/menus", handlers.GetMenus).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
