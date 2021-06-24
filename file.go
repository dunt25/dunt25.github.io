package main

import (
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/image/webp"
)

// readFile will read from a text file to string
func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// writeFile will write from string text into a file
func writeFile(file string, text string) error {
	err := ioutil.WriteFile(file, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(name, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(getFilePath(name))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// convertWebp will convert a webp file into png
func convertWebp(fileName, ext string) error {
	filePath := getFilePath(fileName)

	// Open the file
	file, err := os.Open(filePath + "." + ext)
	if err != nil {
		return err
	}

	// Decode webp file
	img, err := webp.Decode(file)
	if err != nil {
		return err
	}

	// Create new png file
	f, err := os.Create(filePath + ".png")
	if err != nil {
		return err
	}

	// Encode to png
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}
