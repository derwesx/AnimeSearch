package models

import "github.com/lib/pq"

type Anime struct {
	Id          int            `json:"-" json:",omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Name        string         `json:"name" gorm:"not null;default:null"`
	OriginName  string         `json:"origin_name" gorm:"not null;default:null"`
	Description string         `json:"description" gorm:"default:''"`
	Rating      int            `json:"rating" gorm:"not null;default:10"`
	AnimeHash   string         `json:"anime_hash" gorm:"unique;not null;default:null"`
	Episodes    pq.StringArray `json:"episodes,omitempty" gorm:"type:text[];default:ARRAY[]::text[]"`
}
