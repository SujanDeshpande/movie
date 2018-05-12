package files

import (
	"encoding/json"
	"errors"
	"movie/database"
	"movie/utils"
	"time"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

//FileInfo - type of FileInfo
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    uint32    `json:"mode"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
	From    string    `json:"from"`
	To      string    `json:"to"`
}

//Create - CreatefileInfo
func (f FileInfo) Create(fileInfo *FileInfo) error {
	database.DBCon.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.GetConfig().Database.Bucket))
		if bucket == nil {
			berr := errors.New("Bucket not found")
			log.WithError(berr).Error("Bucket not found")
		}
		jvalue, _ := json.Marshal(fileInfo)
		perr := bucket.Put([]byte(fileInfo.Name), jvalue)
		if perr != nil {
			log.WithError(perr).Error("Could not persist movie")
		}
		return nil
	})
	return nil
}

//GetAll - GetAll FileInfo
func (f FileInfo) GetAll() ([]FileInfo, error) {
	var fileInfos []FileInfo
	database.DBCon.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.GetConfig().Database.Bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			jvalue := FileInfo{}
			json.Unmarshal(v, &jvalue)
			fileInfos = append(fileInfos, jvalue)
		}
		return nil
	})
	return fileInfos, nil
}
