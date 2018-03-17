package Utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//ReadFile - reads the File from specified location and returns as a string
func ReadFile(filename string) string {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	str := string(file)
	return str
}

//MakeTimestamp - generates the system timestamp.
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
