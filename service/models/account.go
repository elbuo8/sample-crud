package models

type Account struct {
	ID     string  `gorm:"primary_key:true"`
	Models []Model `gorm:"foreignkey:AccountID;association_foreignkey:ID"`
}
