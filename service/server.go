package main

import (
	"./controllers"
	"./models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func InitDB() *gorm.DB {
	// Type should be postgres as it is coupled at the moment.
	db, err := gorm.Open("postgres", os.Getenv("POSTGRES_CONNECTION_DETAILS"))
	if err != nil {
		log.Fatal(err)
	}
	// Rather than using defer to close the DB, we should listen to OS signals and close accordingly
	return db
}

func RegisterRoutes(r *mux.Router, db *gorm.DB, logger *zap.Logger) {
	m := controllers.NewModel(db, logger)
	r.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})).Methods("GET")
	r.Handle("/models/{id}", m.GetByID()).Methods("GET")
	r.Handle("/models/{id}", m.Post()).Methods("POST")
	r.Handle("/models/{id}", m.Delete()).Methods("DELETE")
	r.Handle("/models", m.Put()).Methods("PUT")
	r.Handle("/models", m.Get()).Methods("GET")
}

func BuildService() http.Handler {
	r := mux.NewRouter()
	db := InitDB()
	db.LogMode(true)
	db.AutoMigrate(&models.Model{}, &models.Account{})
	logger, _ := zap.NewProduction()
	RegisterRoutes(r, db, logger)
	return r
}

func main() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", os.Getenv("PORT")), BuildService()))
}
