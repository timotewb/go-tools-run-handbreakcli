package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Title struct {
	Index     int    `json:index`
	Subtitles string `json:subtitles`
}

func main() {
	fmt.Println("testing")

	var titles []Title

	content, err := ioutil.ReadFile("el.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice to a string
	text := string(content)

	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Loop through each line
	record := false
	subs := false
	for _, line := range lines {
		// Test if we have reached the title details
		if strings.HasPrefix(line, "+ title ") {
			record = true
		}
		if record {
			// Get index
			if strings.Contains(line, "+ index") {
				// Convert the string to an integer
				index, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
					return
				}
			}
			// Test for subtitles
			if strings.Contains(line, "+ subtitle tracks") {
				subs = true
			} else if subs {
				fmt.Println(line)
			}
		}
	}
}
