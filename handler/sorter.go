package Handler

import (
	"errors"
	"movie/Utils"
	"movie/files"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

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
		dbChan := make(chan Files.FileInfo)
		wg.Add(1)
		go func() {
			writeToDB(destChan, dbChan, done)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			moveFile(destChan, done)
			wg.Done()
		}()
		go func() {
			wg.Wait()
			log.Info("Wait over")
		}()
		log.Info("Sent Data")
		w.Write([]byte("\n[]byte\n\n"))
	})
}

func walkFiles(done <-chan struct{}) <-chan Files.FileInfo {
	infoChan := make(chan Files.FileInfo)
	go func() {
		defer close(infoChan)
		defer log.Info("infochan closed")
		filepath.Walk(Utils.GetConfig().Location.Src, func(path string, info os.FileInfo, err error) error {
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
		destYear := Utils.GetConfig().Location.Dest + strconv.Itoa(fileDate.Year())
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
		fileInfo.Create(&fileInfo)
		select {
		case out <- fileInfo:
		case <-done:
			return
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

func createFolder(folderPath string) {
	folder := folderPath
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		merr := os.Mkdir(folder, 0777)
		if merr != nil {

		}
	}
}
