package models

import "github.com/lib/pq"

type User struct {
	Id         int           `gorm:"primary_key;AUTO_INCREMENT"`
	Name       string        `json:"username" gorm:"not null;default:null"`
	Email      string        `json:"email" gorm:"unique"`
	Password   string        `json:"password,omitempty" gorm:"not null;default:null"`
	Salt       []byte        `json:",omitempty" gorm:"not null;default:null"`
	IsAdmin    bool          `json:"is_admin,omitempty" gorm:"default:false"`
	Favourites pq.Int64Array `json:"favourites,omitempty" gorm:"type:integer[]"`
}
