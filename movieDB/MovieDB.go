package MovieDB

import (
	"encoding/json"
	"errors"
	"movie/Structures"
	"movie/Utils"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// MovieStore interface for all movie DB operations
type MovieStore interface {
	AddMovie(Structures.Movie) error
	AddMovies(Structures.MovieInfo) error
	DeleteMovie(string) error
	DeleteMovies([]string) error
	ReadMovieByID(string) error
}

//DB object for database
type DB struct {
	Bolted *bolt.DB
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

// AddMovie Adds a single Movie
func (db *DB) AddMovie(movie Structures.Movie) error {
	db.Bolted.Update(func(tx *bolt.Tx) error {
		movieJSON, jerr := json.Marshal(movie)
		if jerr != nil {
			log.WithError(jerr).WithFields(log.Fields{
				"movieID":   movie.ID,
				"movieName": movie.Name,
			}).Error("Could not create JSON")
			return jerr
		}
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		if bucket == nil {
			berr := errors.New("Bucket not found")
			log.WithError(berr).Error("Bucket not found")
			return berr
		}
		perr := bucket.Put([]byte(movie.ID), movieJSON)
		if perr != nil {
			log.WithError(perr).Error("Could not persist movie")
			return perr
		}
		log.WithFields(log.Fields{
			"movieID":   movie.ID,
			"movieName": movie.Name,
		}).Info("Persisted")
		return nil
	})
	return nil
}

//ReadMovieByID - read movie based on id
func (db *DB) ReadMovieByID(ID string) error {
	db.Bolted.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		if bucket == nil {
			berr := errors.New("Bucket not found")
			log.WithError(berr).Error("Bucket not found")
			return berr
		}
		movieByte := bucket.Get([]byte(ID))
		movie := Structures.Movie{}
		json.Unmarshal(movieByte, &movie)
		log.WithFields(log.Fields{
			"movieID":   movie.ID,
			"movieName": movie.Name,
		}).Info("Retrieved")
		return nil
	})
	return nil
}

//CreateBucket Creates a new Bucket if it does not exist
func CreateBucket(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Utils.GetDatabaseConfig().Bucket))
		if bucket != nil {
			log.WithField("bucket", Utils.GetDatabaseConfig().Bucket).Info("Bucket already exists")
		}
		if bucket == nil {
			_, berr := tx.CreateBucket([]byte(Utils.GetDatabaseConfig().Bucket))
			if berr != nil {
				log.WithError(berr).WithField("bucket", Utils.GetDatabaseConfig().Bucket).Fatal("Unable to create a bucket")
			}
			log.Info("Bucket created sucessfully")
		}
		return nil
	})
}
