package Utils

import (
	"io/ioutil"
	"movie/Structures"
)

// Read - Reads file
func Read(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return filenames
}

// ReadToJSON - Reads file  and converts to a json
func ReadToJSON(path string) Structures.MovieInfo {
	files, _ := ioutil.ReadDir(path)
	movies := []Structures.Movie{}
	for _, f := range files {

		movie := Structures.Movie{}
		movie.Name = f.Name()

		movies = append(movies, movie)

	}
	movieInfo := Structures.MovieInfo{Movies: movies}
	return movieInfo
}
