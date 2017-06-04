package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie/Structures"
	"movie/fileUtils"
	"strings"
)

func main() {
	movies := FileUtils.ReadToJSON("/Volumes/Public/Shared Videos/ENGLISH")
	assignNewName(&movies)
	str, _ := json.Marshal(movies)
	fmt.Println(string(str))
	ioutil.WriteFile("/Users/sdeshpande/Desktop/movie.json", str, 0644)

}

func assignNewName(movies *Structures.MovieInfo) {
	for _, element := range movies.Movies {
		str := strings.Split(element.Name, " ")
		for index, stringer := range str {
			output := fmt.Sprintf("%d) %s",index,stringer )
			fmt.Println(output)

		}
		var input string
		fmt.Scanln(&input)
		fmt.Print(input)
	}
}
