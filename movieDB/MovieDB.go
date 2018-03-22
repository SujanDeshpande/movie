package MovieDB

import (
	"movie/Structures"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// MovieStore interface for all movie DB operations
type MovieStore interface {
	AllMovies() (Structures.MovieInfo, error)
}

//DB object for database
type DB struct {
	Bolted *bolt.DB
}

func init() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// AllMovies return all movies
func (db *DB) AllMovies() (Structures.MovieInfo, error) {
	movies := []Structures.Movie{}

	movie := Structures.Movie{}
	movie.Name = "m1"
	movie.ID = "1"
	movies = append(movies, movie)
	movieInfo := Structures.MovieInfo{Movies: movies}
	log.Info("All movies")
	log.Info(db.Bolted.Info())
	return movieInfo, nil
}
