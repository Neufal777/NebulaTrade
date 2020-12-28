package utils

import (
	"io/ioutil"
	"log"
	"strconv"
)

//WriteFile - writed
func WriteFile(input string, filename string) string {
	err := ioutil.WriteFile(filename, []byte(input), 0)

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

//FloatToString - Convert float to string
func FloatToString(inputnum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputnum, 'f', 12, 64)
}

//StringToFloat - Convert string to float
func StringToFloat(input string) float64 {

	result, _ := strconv.ParseFloat(input, 64)
	return result
}
