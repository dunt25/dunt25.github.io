package main

import (
	"errors"
	"net/url"
	"strings"
)

func downloadImage(idx int, link string) error {
	arrLink := strings.Split(link, " >> ")

	if !validateLink(arrLink[0]) {
		return nil
	}
	totalCount++

	// Get original file name and extension decoded
	path := strings.Split(arrLink[0], "/")
	fullName := path[len(path)-1]
	arrName := strings.Split(fullName, ".")
	ext := arrName[len(arrName)-1]

	if len(arrLink) > 1 {
		// override with custom name
		fullName = arrLink[1] + "." + ext
	}

	name, err := url.QueryUnescape(fullName)
	if err != nil {
		return err
	}

	// Download the file
	err = downloadFile(name, arrLink[0])
	if err != nil {
		return errors.New("Error downloading: " + arrLink[0] + " err: " + err.Error())
	}
	logWithTag("I", idx, "Downloaded: "+arrLink[0])

	// Get raw file name
	fileName := strings.Split(name, ".")

	// Convert if in webp format
	if ext := fileName[len(fileName)-1]; ext == "webp" {
		err = convertWebp(fileName[0], ext)
		if err != nil {
			return errors.New("Error converting: " + arrLink[0] + " err: " + err.Error())
		}
		logWithTag("I", idx, "Converted: "+arrLink[0])
	}

	successCount++
	return nil
}

func getImageLinks(link string) ([]string, []string, error) {
	request := NewRequest()
	request.URL = link
	request.Method = "GET"

	response, body, err := request.doRequest()
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	arr1 := urlRegex.FindAllString(string(body), -1)
	arr2 := rawUrlRegex.FindAllString(string(body), -1)

	return arr1, arr2, nil
}

func readPage(idx int, link string) error {
	logWithTag("P", idx, "Fetching from page: "+link)

	imgLinks, rawImgLinks, err := getImageLinks(link)
	if err != nil {
		return err
	}

	for idx, imgLink := range imgLinks {
		err = downloadImage(idx, imgLink)
		if err != nil {
			logWithTag("P", idx, "Error downloading: "+imgLink+" "+err.Error())
		}
	}

	for idx, imgLink := range rawImgLinks {
		err = downloadImage(idx, "https://"+imgLink)
		if err != nil {
			logWithTag("P", idx, "Error downloading: https://"+imgLink+" "+err.Error())
		}
	}

	return nil
}

func fetchLinks(link string) error {
	imgLinks, rawImgLinks, err := getImageLinks(link)
	if err != nil {
		return err
	}

	var text, text2 string

	text = strings.Join(imgLinks, "\n")
	if len(rawImgLinks) > 0 {
		text2 = "\nhttps://" + strings.Join(rawImgLinks, "\nhttps://")
	}
	outLog.WriteString(text + text2)

	return nil
}
