package models

import "github.com/lib/pq"

type User struct {
	Id        int            `json:"-" gorm:"primary_key;AUTO_INCREMENT"`
	Name      string         `json:"username" gorm:"not null;default:null"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-" gorm:"not null;default:null"`
	Salt      []byte         `json:"-" gorm:"not null;default:null"`
	IsAdmin   bool           `json:"-" gorm:"default:false"`
	Playlists pq.StringArray `json:"playlists,omitempty" gorm:"type:text[];default:ARRAY[]::text[]"`
}
