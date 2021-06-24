package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	mode                     *string
	totalCount, successCount int
	logPrefix                string
	fileLog, outLog          *os.File

	act = "downloading"
)

const (
	folderPath = "images"

	modeImage = "image"
	modePage  = "page"
	modeLink  = "link"
	modeClean = "clean"
)

func main() {
	timeStart := time.Now()
	os.MkdirAll(folderPath, os.ModePerm)

	// Read file
	file, err := readFile("input.txt")
	if err != nil {
		panic(err)
	}

	log, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	fileLog = log
	defer fileLog.Close()

	// Get list of URLs
	urlList := strings.Split(file, "\n")
	reqMode := getMode()

	// Download each URL
	for idx, link := range urlList {
		switch reqMode {
		case modePage:
			if isLoop, newLink, pf, pu := getPageQuery(link); isLoop {
				for i := pf; i <= pu; i++ {
					err = readPage(idx, strings.Replace(newLink, "{page}", strconv.Itoa(i), -1))
				}
			} else {
				err = readPage(idx, link)
			}
			logPrefix = "P"

		case modeImage:
			err = downloadImage(idx, link)
			logPrefix = "I"

		case modeLink:
			outLog, err = os.Create("output.txt")
			if err != nil {
				panic(err)
			}
			defer outLog.Close()

			if isLoop, newLink, pf, pu := getPageQuery(link); isLoop {
				for i := pf; i <= pu; i++ {
					err = fetchLinks(strings.Replace(newLink, "{page}", strconv.Itoa(i), -1))

				}
			} else {
				err = fetchLinks(link)
			}
			logPrefix = "I"
			act = "fetching links"

		case modeClean:
			os.RemoveAll(folderPath)
			act = "deleting"

		default:
			err = errors.New("Undefined mode")
			logPrefix = "U"
		}

		if err != nil {
			logWithTag(logPrefix, idx, err.Error())

			if logPrefix == "U" {
				break
			}
		}
	}

	str := "image"
	if successCount > 1 {
		str = "images"
	}

	var counterStr string
	if successCount > 0 {
		counterStr = fmt.Sprintf(" %d/%d %s", successCount, totalCount, str)
	}

	fmt.Printf("Finished %s%s in %v\n", act, counterStr, time.Now().Sub(timeStart))
}
