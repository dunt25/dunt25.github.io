package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	urlRegex    = regexp.MustCompile(`(http(s?):)([/|.|\w|\s|%|-])*\.(?:jpg|jpeg|webp|png)`)
	rawUrlRegex = regexp.MustCompile(`(www)([/|.|\w|\s|-])*\.(?:jpg|jpeg|webp|png)`)
)

func logWithTag(prefix string, idx int, msg string) {
	text := fmt.Sprintf("[%s%d] %s\n", prefix, idx+1, msg)
	fmt.Print(text)
	fileLog.WriteString(text)
}

func validateLink(link string) bool {
	return urlRegex.MatchString(link)
}

func getFilePath(name string) string {
	return folderPath + "/" + name
}

func getMode() string {
	if mode != nil {
		return *mode
	}

	var m string
	if len(os.Args) > 1 {
		m = os.Args[1][1:]
	}
	mode = &m

	return *mode
}

func getPageQuery(link string) (bool, string, int, int) {
	arrLink := strings.Split(link, " ")
	if len(arrLink) < 2 {
		return false, arrLink[0], 0, 0
	}

	var (
		pageStart, pageEnd int
		err                error
	)

	pageStart, err = strconv.Atoi(arrLink[1])
	if err != nil {
		return false, arrLink[0], 0, 0
	}

	if len(arrLink) > 2 {
		pageEnd, err = strconv.Atoi(arrLink[2])
		if err != nil || pageEnd < pageStart {
			pageEnd = pageStart
		}
	}

	return true, arrLink[0], pageStart, pageEnd
}
