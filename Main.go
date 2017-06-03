package main

import (
	"movie/fileUtils"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func main() {
	movies := FileUtils.ReadToJSON("/Volumes/Public/Shared Videos/ENGLISH")
	str, _ := json.Marshal(movies)
	fmt.Println(string(str))
	ioutil.WriteFile("/Users/sdeshpande/Desktop/movie.json", str, 0644)



}
