package controllers

import (
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"net/http"
)

// In reality I'd make models more self contained.
// Basically it would handle the DB connection and the controllers/middlewares would simply work with the models.
// But for simplicity went with this initial route.

type Model struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewModel(db *gorm.DB, logger *zap.Logger) *Model {
	m := Model{}
	m.db = db
	m.logger = logger
	return &m
}

func (c *Model) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		m := models.Model{}
		results := []models.Model{}
		// Figure out how to dynamically set this
		if name := params.Get("name"); name != "" {
			m.Name = name
		}

		if accountId := params.Get("accountId"); accountId != "" {
			m.AccountID = accountId
		}

		query := c.db.Where(&m)
		if sort := params.Get("sortBy"); sort != "" {
			query = query.Order(fmt.Sprintf("%s desc", sort))
		}

		if total := params.Get("total"); total != "" {
			query = query.Limit(total)
		}

		if err := query.Find(&results).Error; err != nil {
			c.logger.Error(err.Error())
			http.Error(w, "Something blew up", 500)
			return
		}

		json.NewEncoder(w).Encode(results)
	})
}

func (c *Model) GetByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		m := models.Model{ID: id}
		if err := c.db.Find(&m).Error; err != nil {
			// Check if not found or ISR
			if gorm.IsRecordNotFoundError(err) {
				http.NotFound(w, r)
			} else {
				c.logger.Error(err.Error())
				http.Error(w, "Something blew up", 500)
			}
			return
		}

		json.NewEncoder(w).Encode(m)
	})
}

func (c *Model) Put() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := models.Model{}
		// I'd move these validations into middlewares
		if r.Body == nil {
			http.Error(w, "Empty payload", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, "Invalid payload", 400)
			return
		}

		m.ID = ksuid.New().String()

		// Should validate that it was created by the proper owner
		// But for now, put them all under the same owner
		m.AccountID = "1AMMDguYEpokKpOpXjRIFBezvVd"

		if err = c.db.Create(&m).Error; err != nil {
			// depending on the reason, log it
			http.Error(w, "Failed to create model", 400)
			return
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(m)
	})
}

func (c *Model) Post() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		m := models.Model{ID: id}
		if err := c.db.Find(&m).Error; err != nil {
			// Check if not found or ISR
			if gorm.IsRecordNotFoundError(err) {
				http.NotFound(w, r)
			} else {
				c.logger.Error(err.Error())
				http.Error(w, "Something blew up", 500)
			}
			return
		}

		if r.Body == nil {
			http.Error(w, "Empty payload", 400)
			return
		}
		payload := models.Model{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid payload", 400)
			return
		}
		// stricter validations (avoid the ID or the owner to be changed)
		if err = mergo.Merge(&m, payload, mergo.WithOverride); err != nil {
			c.logger.Error(err.Error())
			http.Error(w, "Something blew up", 500)
		}
		c.db.Save(&m)
		json.NewEncoder(w).Encode(m)
	})
}

func (c *Model) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.NotFound(w, r)
			return
		}
		m := models.Model{ID: id}
		if err := c.db.Find(&m).Error; err != nil {
			// Check if not found or ISR
			if gorm.IsRecordNotFoundError(err) {
				http.NotFound(w, r)
			} else {
				c.logger.Error(err.Error())
				http.Error(w, "Something blew up", 500)
			}
			return
		}
		if err := c.db.Delete(&m).Error; err != nil {
			c.logger.Error(err.Error())
			http.Error(w, "Something blew up", 500)
			return
		}
		w.WriteHeader(204)
	})
}
