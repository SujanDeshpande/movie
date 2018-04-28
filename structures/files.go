package Structures

import "time"

type FileInfo struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Mode     uint32    `json:"mode"`
	ModTime  time.Time `json:"modTime"`
	IsDir    bool      `json:"isDir"`
	Location string    `json:"location"`
}
