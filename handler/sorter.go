package Handler

import (
	"encoding/json"
	"errors"
	"movie/MovieDB"
	"movie/Utils"
	"movie/files"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

//SortHandler is Default Sort API Handler
func SortHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		infoChan := walkFiles(done)
		var wg sync.WaitGroup
		wg.Add(1)
		destChan := make(chan Files.FileInfo)
		go func() {
			createDestination(infoChan, destChan, done)
			wg.Done()

		}()
		go func() {
			wg.Wait()
			log.Info("Wait over")
		}()
		dbChan := make(chan Files.FileInfo)
		wg.Add(1)
		go func() {
			writeToDB(destChan, dbChan, done)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			moveFile(dbChan, done)
			wg.Done()
		}()

		log.Info("Sent Data")
		w.Write([]byte("\n[]byte\n\n"))

		// destChan := createDestination(infoChan, destFolder)
		// writeChan := writeToDB(destChan)
		// moveChan := moveFile(writeChan, srcFolder)
		//
		// for fileInfoq := range moveChan {
		// 	log.Info(fileInfoq.ModTime)
		// }
	})
}

func walkFiles(done <-chan struct{}) <-chan Files.FileInfo {
	infoChan := make(chan Files.FileInfo)
	go func() {
		defer close(infoChan)
		defer log.Info("infochan closed")
		filepath.Walk(Utils.GetConfig().Files.Src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileInfo := Files.FileInfo{Name: info.Name(),
				Mode:    uint32(info.Mode()),
				ModTime: info.ModTime(),
				Size:    info.Size(),
				IsDir:   info.IsDir(),
				From:    path,
			}
			select {
			case infoChan <- fileInfo:
			case <-done:
				return errors.New("Walk Canceled")
			}
			return nil
		})
	}()
	return infoChan
}

func createDestination(in <-chan Files.FileInfo, out chan<- Files.FileInfo, done <-chan struct{}) {
	defer close(out)
	defer log.Info("createDestination closed")
	for fileInfo := range in {
		fileDate := fileInfo.ModTime
		destYear := Utils.GetConfig().Files.Dest + strconv.Itoa(fileDate.Year())
		createFolder(destYear)
		destMonth := destYear + "/" + fileDate.Month().String()
		createFolder(destMonth)
		fileInfo.To = destMonth
		select {
		case out <- fileInfo:
		case <-done:
			return
		}
	}
}

func writeToDB(in <-chan Files.FileInfo, out chan<- Files.FileInfo, done <-chan struct{}) {
	defer close(out)
	defer log.Info("writeToDB closed")
	for fileInfo := range in {
		fileInfo.CreateFileInfo(&fileInfo)
		select {
		case out <- fileInfo:
		case <-done:
		}
	}
}

func moveFile(in <-chan Files.FileInfo, done chan<- struct{}) {
	defer log.Info("moveFile closed")
	for fileInfo := range in {
		log.Info(fileInfo.IsDir)
		if !fileInfo.IsDir {
			fromFile := fileInfo.From
			toFile := fileInfo.To + "/" + fileInfo.Name
			os.Rename(fromFile, toFile)
			log.Info(fromFile + " Moved to " + toFile)
		} else {
			log.Info(fileInfo.Name + "Not Moved")
		}
	}
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

func createFolder(folderPath string) {
	folder := folderPath
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		merr := os.Mkdir(folder, 0777)
		if merr != nil {

		}
	}
}
