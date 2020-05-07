package models

type Album  struct {
	AlbumName    string 	`json:"albumName"`
	Images     []string		`json:"imageName"`
}