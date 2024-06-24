package models

type Anime struct {
	Id         int    `json:",omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string `json:"name" gorm:"not null;default:null"`
	OriginName string `json:"origin_name" gorm:"not null;default:null"`
	AnimeHash  string `json:"anime_hash" gorm:"unique;not null;default:null"`
}
