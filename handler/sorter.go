package Handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"movie/MovieDB"
	"movie/Utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

//SortHandler is Default Sort API Handler
func SortHandler(db MovieDB.DB) http.Handler {
	config := Utils.GetConfig()
	srcFolder := config.Files.Src
	destFolder := config.Files.Dest
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infoChan := getFileInfo(srcFolder)
		destChan := createDestination(infoChan, destFolder)
		writeChan := writeToDB(destChan, db)
		moveChan := moveFile(writeChan, srcFolder)

		for filerq := range moveChan {
			log.Info(filerq.ModTime)
			filers := getFileInfoFromDB(db)
			for _, filere := range filers {
				log.Info("sss")
				log.Info(filere.ModTime)
			}
		}

	})
}

func writeToDB(inChan <-chan filer, db MovieDB.DB) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for filer := range inChan {
			addFileInfoToDB(filer, db)
			outChan <- filer
		}
	}()
	return outChan
}

type filer struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Mode        uint32    `json:"mode"`
	ModTime     time.Time `json:"modTime"`
	IsDir       bool      `json:"isDir"`
	Destination string    `json:"destination"`
}

// AddFileInfoToDB Adds a single Movie
func addFileInfoToDB(filer filer, db MovieDB.DB) {
	db.Bolted.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		if bucket == nil {
			berr := errors.New("Bucket not found")
			log.WithError(berr).Error("Bucket not found")
		}
		jvalue, _ := json.Marshal(filer)
		perr := bucket.Put([]byte(filer.Name), jvalue)
		if perr != nil {
			log.WithError(perr).Error("Could not persist movie")
		}
		return nil
	})
}

//getFileInfoFromDB - read movie based on id
func getFileInfoFromDB(db MovieDB.DB) []filer {
	var filers []filer
	db.Bolted.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			jvalue := filer{}
			json.Unmarshal(v, &jvalue)
			filers = append(filers, jvalue)
		}
		return nil
	})
	return filers
}

func getFileInfo(srcFolder string) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		files, _ := ioutil.ReadDir(srcFolder)
		for _, f := range files {
			filerw := filer{Name: f.Name(),
				Mode:    uint32(f.Mode()),
				ModTime: f.ModTime(),
				Size:    f.Size(),
				IsDir:   f.IsDir(),
			}

			outChan <- filerw
		}

	}()
	return outChan
}

func createDestination(inChan <-chan filer, dest string) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for info := range inChan {
			datee := info.ModTime
			destYear := dest + strconv.Itoa(datee.Year())
			createFolder(destYear)
			destMonth := destYear + "/" + datee.Month().String()
			createFolder(destMonth)
			info.Name = destMonth
			outChan <- info
		}
	}()
	return outChan
}
func moveFile(inChan <-chan filer, src string) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for filer := range inChan {
			fromFile := src + "/" + filer.Name
			toFile := filer.Destination + "/" + filer.Name
			os.Rename(fromFile, toFile)
			outChan <- filer
		}
	}()
	return outChan
}

func createFolder(folderPath string) {
	folder := folderPath
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		merr := os.Mkdir(folder, 0777)
		if merr != nil {

		}
	}
}
