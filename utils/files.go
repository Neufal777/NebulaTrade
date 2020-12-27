package utils

import (
	"fmt"
	"io/ioutil"
	"log"
)

//WriteFile - writed
func WriteFile(input float64, filename string) string {
	err := ioutil.WriteFile(filename, []byte(fmt.Sprintf("%f", input)), 0)

	if err != nil {
		log.Fatal(err)
	}

	return "Writed correctly :)"
}

//ReadFile - Read file
func ReadFile(filename string) string {

	fileContent, err := ioutil.ReadFile("core/status.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileContentString := string(fileContent)

	return fileContentString
}
