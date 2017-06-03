package FileUtils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie/Structures"
)

func Read(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func ReadToJSON(path string) {
	files, _ := ioutil.ReadDir(path)
	movies := []Structures.Movie{}
	for _, f := range files {

		movie := Structures.Movie{}
		movie.Name = f.Name()
		movie.IsFolder = f.IsDir()
		if !f.IsDir() {
			movie.FileSize = f.Size()
		}
		movies = append(movies, movie)

	}
	movieInfo := Structures.MovieInfo{movies}
	str, _ := json.Marshal(movieInfo)
	fmt.Println(string(str))
}
