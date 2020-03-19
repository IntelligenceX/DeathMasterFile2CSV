/*
File Name:  Converter.go
Copyright:  2020 Kleissner Investments s.r.o.
Author:     Peter Kleissner

Use: Converter [Input File] [Output File]
*/

package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Converter [Input File] [Output File]")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	convertDMF(inputFile, outputFile)
}

const readAtOnce = 100 * 1024 * 1024 // 100 MB

// splitLastNL splits the data at the last new-line, if available
func splitLastNL(data []byte) (first, second []byte) {
	if index := bytes.LastIndex(data, []byte{'\n'}); index > 0 {
		first = data[:index]
		second = data[index:]
		return
	}

	return first, nil
}
