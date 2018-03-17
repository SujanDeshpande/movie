package FileUtils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie/Structures"
	"regexp"
	"strings"
)

// Something - Something wasted
func Something() {
	movies := ReadToJSON("/Volumes/Public/Shared Videos/ENGLISH")
	assignNewName(&movies)
	str, _ := json.Marshal(movies)
	fmt.Println(string(str))
	ioutil.WriteFile("/Users/sdeshpande/Desktop/movie.json", str, 0644)

}

func assignNewName(movies *Structures.MovieInfo) {
	for _, element := range movies.Movies {
		var str []string
		re := regexp.MustCompile("(.*\\.)+.*")
		if len(re.FindString(element.Name)) > 0 {
			str = strings.Split(element.Name, ".")
		} else {
			str = strings.Split(element.Name, " ")
		}
		for index, stringer := range str {
			output := fmt.Sprintf("%d) %s", index, stringer)
			fmt.Println(output)

		}
		var inputStr string
		fmt.Scanln(&inputStr)
		input := strings.Split(inputStr, " ")

		convertStringToIntArray(input)
	}
}

func convertStringToIntArray(strings []string) {
	for _, i := range strings {
		fmt.Println("sdsads")
		fmt.Println(i)
	}
}
