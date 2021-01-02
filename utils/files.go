package utils

import (
	"fmt"
	"strconv"
)

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

//AnyTypeToString -
func AnyTypeToString(input interface{}) string {
	return fmt.Sprintf("%v", input)
}
