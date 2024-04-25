package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"os"
)

func GenerateReadme(filename string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://masagu.dev/index.xml")
	if err != nil {
		log.Fatalf("error getting feed: %v", err)
	}

	blogItem := feed.Items[0]

	pre, err := ReadFileAsString("./PREREADME.md")
	blog := "- Latest blog post :page_facing_up: [" + blogItem.Title + "](" + blogItem.Link + ")"
	post, err := ReadFileAsString("./POSTREADME.md")
	data := fmt.Sprintf("%s\n%s\n%s", pre, blog, post)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func ReadFileAsString(filename string) (string, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading file: %s %v", filename, err)
	}
	return string(dat), nil
}

func main() {
	GenerateReadme("../README.md")
}
