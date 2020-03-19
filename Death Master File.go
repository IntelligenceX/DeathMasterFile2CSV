/*
File Name:  SSN Death Master File.go
Copyright:  2020 Kleissner Investments s.r.o.
Author:     Peter Kleissner

Converts the proprietary format of the SSN Death Master File to CSV.
*/

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
SSN Death Master File Documentation:
* https://en.wikipedia.org/wiki/Death_Master_File
* https://dmf.ntis.gov/recordlayout.pdf

POSITION 1 DESCRIPTION -BLANK, OR A (ADD), C (CHANGE), OR D (DELETE) *** FIELD SIZE - 1
POSITION 2-10 DESCRIPTION -SOCIAL SECURITY NUMBER FIELD SIZE - 9
POSITION 11-30 DESCRIPTION - LAST NAME FIELD SIZE – 20 *
POSITION 31-34 DESCRIPTION – NAME SUFFIX FIELD SIZE – 4 *
POSITION 35-49 DESCRIPTION – FIRST NAME FIELD SIZE – 15 *
POSITION 50-64 DESCRIPTION – MIDDLE NAME FIELD SIZE – 15 *
POSITION 65 DESCRIPTION - V OR P CODE (VERIFIED OR PROOF CODE) *** FIELD SIZE – 1 *
POSITION 66-73 DESCRIPTION – DATE OF DEATH (MM,DD,CC,YY) FIELD SIZE – 8 *
POSITION 74-81 DESCRIPTION – DATE OF BIRTH (MM,DD,CC,YY) FIELD SIZE – 8 *
POSITION 82-83 DESCRIPTION – BLANKS ***** **
POSITION 84-88 DESCRIPTION – BLANKS ***** **
POSITION 89-93 DESCRIPTION – BLANKS ***** **
POSITION 94-100 DESCRIPTION – BLANKS
*/

var csvHeaderDMF = []string{"Type", "Social Security Number", "Last Name", "Name Suffix", "First Name", "Middle Name", "Verified", "Date of Death", "Date of Birth", "Blank 1", "Blank2", "Blank 3", "Blank 4"}

func convertDMF(inputFile, outputFile string) {
	// open input file, create output file
	fileI, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		return
	}
	defer fileI.Close()

	fileO, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		return
	}
	defer fileO.Close()

	// output new CSV header
	outputCSV := csv.NewWriter(fileO)
	outputCSV.Write(csvHeaderDMF)
	outputCSV.Flush()

	var pre []byte

	// read & process loop
	for {
		bufferRead := make([]byte, readAtOnce)
		n1, err := fileI.Read(bufferRead)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading from input file: %s\n", err)
			return
		}

		first, second := splitLastNL(bufferRead[:n1])
		first = append(pre, first...)
		pre = second

		fmt.Printf("Read %d bytes, process %d, new line after is %d bytes\n", n1, len(first), len(second))

		count := convertDMF2CSV(first, outputCSV)

		fmt.Printf("Wrote %d records\n", count)
	}

	convertDMF2CSV(pre, outputCSV)

	fmt.Printf("All done!\n")
}

func convertDMF2CSV(input []byte, output *csv.Writer) (countRecords int) {
	lines := bytes.Split(input, []byte("\n"))

	for n := range lines {
		lineA := string(lines[n])

		if len(lineA) == 0 {
			continue
		} else if len(lineA) != 100 {
			fmt.Printf("Skip line, mismatch %d length (expected 100)\n", len(lineA))
			continue
		}

		// extract the fields & transform
		fieldType := strings.TrimSpace(lineA[0:1])          // POSITION 1 DESCRIPTION -BLANK, OR A (ADD), C (CHANGE), OR D (DELETE) *** FIELD SIZE - 1
		fieldSSN := strings.TrimSpace(lineA[1:10])          // POSITION 2-10 DESCRIPTION -SOCIAL SECURITY NUMBER FIELD SIZE - 9
		fieldLastName := strings.TrimSpace(lineA[10:30])    // POSITION 11-30 DESCRIPTION - LAST NAME FIELD SIZE – 20 *
		fieldNameSuffix := strings.TrimSpace(lineA[30:34])  // POSITION 31-34 DESCRIPTION – NAME SUFFIX FIELD SIZE – 4 *
		fieldFirstName := strings.TrimSpace(lineA[34:49])   // POSITION 35-49 DESCRIPTION – FIRST NAME FIELD SIZE – 15 *
		fieldMiddleName := strings.TrimSpace(lineA[49:64])  // POSITION 50-64 DESCRIPTION – MIDDLE NAME FIELD SIZE – 15 *
		fieldVerified := strings.TrimSpace(lineA[64:65])    // POSITION 65 DESCRIPTION - V OR P CODE (VERIFIED OR PROOF CODE) *** FIELD SIZE – 1 *
		fieldDateOfDeath := strings.TrimSpace(lineA[65:73]) // POSITION 66-73 DESCRIPTION – DATE OF DEATH (MM,DD,CC,YY) FIELD SIZE – 8 *
		fieldDateofBirth := strings.TrimSpace(lineA[73:81]) // POSITION 74-81 DESCRIPTION – DATE OF BIRTH (MM,DD,CC,YY) FIELD SIZE – 8 *
		fieldBlank1 := strings.TrimSpace(lineA[81:83])      // POSITION 82-83 DESCRIPTION – BLANKS ***** **
		fieldBlank2 := strings.TrimSpace(lineA[83:88])      // POSITION 84-88 DESCRIPTION – BLANKS ***** **
		fieldBlank3 := strings.TrimSpace(lineA[88:93])      // POSITION 89-93 DESCRIPTION – BLANKS ***** **
		fieldBlank4 := strings.TrimSpace(lineA[93:100])     // POSITION 94-100 DESCRIPTION – BLANKS

		switch fieldType {
		case "A":
			fieldType = "Add"
		case "C":
			fieldType = "Change"
		case "D":
			fieldType = "Delete"
		}

		switch fieldVerified {
		case "V":
			fieldVerified = "Verified"
		case "P":
			fieldVerified = "Proof Code"
		}

		fieldDateOfDeath = convertDMFDate(fieldDateOfDeath)
		fieldDateofBirth = convertDMFDate(fieldDateofBirth)

		fieldLastName = camelTitle(fieldLastName)
		fieldNameSuffix = camelTitle(fieldNameSuffix)
		fieldFirstName = camelTitle(fieldFirstName)
		fieldMiddleName = camelTitle(fieldMiddleName)

		err := output.Write([]string{fieldType, fieldSSN, fieldLastName, fieldNameSuffix, fieldFirstName, fieldMiddleName, fieldVerified, fieldDateOfDeath, fieldDateofBirth, fieldBlank1, fieldBlank2, fieldBlank3, fieldBlank4})

		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
			break
		}

		countRecords++
	}

	output.Flush()

	return
}

func convertDMFDate(encoded string) (regular string) {
	if len(encoded) != 8 {
		return encoded
	}

	// "MMDDCCYY" -> "YYYY-MM-DD"
	year := encoded[4:8]
	month := encoded[0:2]
	day := encoded[2:4]

	return year + "-" + month + "-" + day
}

func camelTitle(title string) string {
	return strings.Title(strings.ToLower(title))
}
