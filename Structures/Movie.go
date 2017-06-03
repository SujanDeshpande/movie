package Structures

type MovieInfo struct {
	Movies []Movie `json:"movies"`
}

type Movie struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsFolder bool   `json:"isFolder"`
	FileSize int64    `json:"fileSize"`
}
