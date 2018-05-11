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

//Create - CreatefileInfo
func (f FileInfo) Create(fileInfo *FileInfo) error {
	db, err := bolt.Open(Utils.GetConfig().Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetConfig().Database.Bucket))
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
	db, err := bolt.Open(Utils.GetConfig().Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetConfig().Database.Bucket))
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

//CreateBucket Creates a new Bucket for Files if it does not exist
func CreateBucket() {
	db, err := bolt.Open(Utils.GetConfig().Database.Name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetConfig().Database.Bucket))
		if bucket != nil {
			log.WithField("bucket", Utils.GetConfig().Database.Bucket).Info("Bucket already exists")
		}
		if bucket == nil {
			_, berr := tx.CreateBucket([]byte(Utils.GetConfig().Database.Bucket))
			if berr != nil {
				log.WithError(berr).WithField("bucket", Utils.GetConfig().Database.Bucket).Fatal("Unable to create a bucket")
			}
			log.Info("Bucket created sucessfully")
		}
		return nil
	})
}
