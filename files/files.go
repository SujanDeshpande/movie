package Files

import (
	"encoding/json"
	"errors"
	"movie/Utils"
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

//CreateFileInfo - CreatefileInfo
func (f FileInfo) CreateFileInfo(fileInfo *FileInfo) error {
	config := Utils.GetConfig()
	db, err := bolt.Open(config.Files.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.Files.Database.Bucket))
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

//GetAllFileInfo - GetAll FileInfo
func GetAllFileInfo() ([]FileInfo, error) {
	var fileInfos []FileInfo
	config := Utils.GetConfig()
	db, err := bolt.Open(config.Files.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.Files.Database.Bucket))
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

//CreateFilesBucket Creates a new Bucket for Files if it does not exist
func CreateFilesBucket() {
	config := Utils.GetConfig()
	db, err := bolt.Open(config.Files.Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.Files.Database.Bucket))
		if bucket != nil {
			log.WithField("bucket", config.Files.Database.Bucket).Info("Bucket already exists")
		}
		if bucket == nil {
			_, berr := tx.CreateBucket([]byte(config.Files.Database.Bucket))
			if berr != nil {
				log.WithError(berr).WithField("bucket", config.Files.Database.Bucket).Fatal("Unable to create a bucket")
			}
			log.Info("Bucket created sucessfully")
		}
		return nil
	})
}
