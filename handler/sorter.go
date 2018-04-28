package Handler

import (
	"io/ioutil"
	"movie/MovieDB"
	"movie/Utils"
	"net/http"
	"os"
	"strconv"

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
		for filer := range moveChan {
			log.Info("processed" + filer.info.Name())
		}
	})
}

type filer struct {
	info        os.FileInfo
	destination string
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
