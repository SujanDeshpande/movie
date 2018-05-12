package handle

import (
	"errors"
	"movie/files"
	"movie/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

func buildPipe() {
	done := make(chan struct{})
	infoChan := walkfiles(done)
	var wg sync.WaitGroup
	wg.Add(1)
	destChan := make(chan files.FileInfo)
	go func() {
		createDestination(infoChan, destChan, done)
		wg.Done()

	}()
	dbChan := make(chan files.FileInfo)
	wg.Add(1)
	go func() {
		writeToDB(destChan, dbChan, done)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		moveFile(dbChan)
		wg.Done()
	}()
	go func() {
		log.Info("Waiting")
		wg.Wait()
		log.Info("Processed")

	}()
}
func walkfiles(done chan struct{}) <-chan files.FileInfo {
	infoChan := make(chan files.FileInfo)
	go func() {
		defer close(infoChan)
		filepath.Walk(utils.GetConfig().Location.Src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileInfo := files.FileInfo{Name: info.Name(),
				Mode:    uint32(info.Mode()),
				ModTime: info.ModTime(),
				Size:    info.Size(),
				IsDir:   info.IsDir(),
				From:    path,
			}
			if fileInfo.IsDir {
				log.WithFields(log.Fields{"directory": fileInfo.Name}).Debug("Parsing")
				return nil
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

func createDestination(in <-chan files.FileInfo, out chan<- files.FileInfo, done <-chan struct{}) {
	defer close(out)
	for fileInfo := range in {
		fileDate := fileInfo.ModTime
		destYear := utils.GetConfig().Location.Dest + strconv.Itoa(fileDate.Year())
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
	log.Info("createDestination")
}

func writeToDB(in <-chan files.FileInfo, out chan<- files.FileInfo, done <-chan struct{}) {
	defer close(out)
	for fileInfo := range in {
		fileInfo.Create(&fileInfo)
		select {
		case out <- fileInfo:
		case <-done:
			return
		}
	}
	log.Info("writeToDB")
}

func moveFile(in <-chan files.FileInfo) {
	for fileInfo := range in {
		if !fileInfo.IsDir {
			fromFile := fileInfo.From
			toFile := fileInfo.To + "/" + fileInfo.Name
			os.Rename(fromFile, toFile)
			log.WithFields(log.Fields{
				"from": fromFile,
				"to":   toFile,
			}).Info("Moved")
		} else {
			log.WithFields(log.Fields{"directory": fileInfo.Name}).Info("Not Moved")
		}
	}
	log.Info("moveFile")
}

func createFolder(folderPath string) {
	folder := folderPath
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		merr := os.Mkdir(folder, 0777)
		if merr != nil {

		}
	}
}
