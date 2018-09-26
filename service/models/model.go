package models

import (
	"github.com/lib/pq"
	"time"
)

// Coupled to psql, could be managed a bit more agnostically  but for simplicity, decided to use psql.
// Missing key things such as "required" etc

type Model struct {
	ID              string         `gorm:"primary_key:true"`
	Name            string         `json:"name,omitempty"`
	Accuracy        float64        `json:"accuracy"`
	AccountID       string         `json:"-"`
	FeatureNames    pq.StringArray `gorm:"type:varchar(64)[]" json:"feature_names,omitempty"`
	HyperParameters pq.StringArray `gorm:"type:varchar(64)[]" json:"hyper_parameters,omitempty"`
	TrainStartTime  *time.Time     `json:"train_start_time,omitempty"`
	TrainStopTime   *time.Time     `json:"train_stop_time,omitempty"`
}
