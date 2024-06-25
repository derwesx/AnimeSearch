package models

import "github.com/lib/pq"

type Playlist struct {
	Id        int            `json:"-" gorm:"primary_key;AUTO_INCREMENT"`
	Name      string         `json:"name" gorm:"not null;default:null"`
	Hash      string         `json:"hash" gorm:"unique"`
	Author    string         `json:"author,omitempty" gorm:"not null;default:null"`
	IsPrivate bool           `json:"-" gorm:"not null;default:true"`
	Animes    pq.StringArray `json:"animes,omitempty" gorm:"type:text[];default:ARRAY[]::text[]"`
}
