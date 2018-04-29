package Handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"movie/MovieDB"
	"movie/Utils"
	"movie/files"
	"net/http"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

//SortHandler is Default Sort API Handler
func SortHandler() http.Handler {
	config := Utils.GetConfig()
	srcFolder := config.Files.Src
	destFolder := config.Files.Dest
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infoChan := getFileInfo(srcFolder)
		destChan := createDestination(infoChan, destFolder)
		writeChan := writeToDB(destChan)
		moveChan := moveFile(writeChan, srcFolder)

		for fileInfoq := range moveChan {
			log.Info(fileInfoq.ModTime)
		}
	})
}

func writeToDB(in <-chan Files.FileInfo) <-chan Files.FileInfo {
	out := make(chan Files.FileInfo)
	go func() {
		defer close(out)
		for fileInfo := range in {
			fileInfo.CreateFileInfo(&fileInfo)
			out <- fileInfo
		}
	}()
	return out
}

// AddFileInfoToDB Adds a single Movie
func addFileInfoToDB(fileInfo Files.FileInfo, db MovieDB.DB) {
	db.Bolted.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
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
}

//getFileInfoFromDB - read movie based on id
func getFileInfoFromDB(db MovieDB.DB) []Files.FileInfo {
	var fileInfos []Files.FileInfo
	db.Bolted.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			jvalue := Files.FileInfo{}
			json.Unmarshal(v, &jvalue)
			fileInfos = append(fileInfos, jvalue)
		}
		return nil
	})
	return fileInfos
}

func getFileInfo(src string) <-chan Files.FileInfo {
	out := make(chan Files.FileInfo)
	go func() {
		defer close(out)
		files, _ := ioutil.ReadDir(src)
		for _, file := range files {
			fileInfo := Files.FileInfo{Name: file.Name(),
				Mode:    uint32(file.Mode()),
				ModTime: file.ModTime(),
				Size:    file.Size(),
				IsDir:   file.IsDir(),
			}

			out <- fileInfo
		}

	}()
	return out
}

func createDestination(in <-chan Files.FileInfo, dest string) <-chan Files.FileInfo {
	out := make(chan Files.FileInfo)
	go func() {
		defer close(out)
		for fileInfo := range in {
			fileDate := fileInfo.ModTime
			destYear := dest + strconv.Itoa(fileDate.Year())
			createFolder(destYear)
			destMonth := destYear + "/" + fileDate.Month().String()
			createFolder(destMonth)
			fileInfo.Location = destMonth
			out <- fileInfo
		}
	}()
	return out
}
func moveFile(in <-chan Files.FileInfo, src string) <-chan Files.FileInfo {
	out := make(chan Files.FileInfo)
	go func() {
		defer close(out)
		for fileInfo := range in {
			fromFile := src + fileInfo.Name
			toFile := fileInfo.Location + "/" + fileInfo.Name
			os.Rename(fromFile, toFile)
			log.Info(fromFile + " Moved to " + toFile)
			out <- fileInfo
		}
	}()
	return out
}

func createFolder(folderPath string) {
	folder := folderPath
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		merr := os.Mkdir(folder, 0777)
		if merr != nil {

		}
	}
}
