package models

type Review struct {
	Id         int       `json:"id"`
	Rating     int		 `json:"rating"`
	Comment    string    `json:"comment"`
}