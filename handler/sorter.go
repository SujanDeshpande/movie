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
		moveChan := moveFile(destChan, srcFolder)
		writeChan := writeToDB(moveChan, db)

		for filer := range writeChan {
			log.Info("processed" + filer.info.Name())
			filers := getFileInfoFromDB(db)
			for i, _ := range filers {
				log.Info("Retrieved" + strconv.Itoa(i))
			}

		}

	})
}

func writeToDB(inChan <-chan filer, db MovieDB.DB) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for filer := range inChan {
			log.Info("writing")
			addFileInfoToDB(filer, db)
			outChan <- filer
		}
	}()
	return outChan
}

type filer struct {
	info        os.FileInfo
	destination string
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
		perr := bucket.Put([]byte(filer.info.Name()), jvalue)
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

func getFileInfo(srcFolder string) <-chan os.FileInfo {
	outChan := make(chan os.FileInfo)
	go func() {
		defer close(outChan)
		files, _ := ioutil.ReadDir(srcFolder)
		for _, f := range files {
			outChan <- f
		}

	}()
	return outChan
}

func createDestination(inChan <-chan os.FileInfo, dest string) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for info := range inChan {
			datee := info.ModTime()
			destYear := dest + strconv.Itoa(datee.Year())
			createFolder(destYear)
			destMonth := destYear + "/" + datee.Month().String()
			createFolder(destMonth)
			outChan <- filer{info, destMonth}
		}
	}()
	return outChan
}
func moveFile(inChan <-chan filer, src string) <-chan filer {
	outChan := make(chan filer)
	go func() {
		defer close(outChan)
		for filer := range inChan {
			fromFile := src + "/" + filer.info.Name()
			toFile := filer.destination + "/" + filer.info.Name()
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
