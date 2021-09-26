package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const dest = ".nvimrc"

//go:embed templates/*
var templates embed.FS

func getLanguages() ([]string, error) {
	var languages = make([]string, 0)
	langTemplates, err := templates.ReadDir("templates")
	if err != nil {
		return nil, errors.New("Unable to read templates")
	}
	for _, template := range langTemplates {
		languages = append(languages, string(template.Name()))
	}
	return languages, nil
}

func templateExists(language string) bool {
	_, err := templates.Open("templates/" + language)
	return err == nil
}

func copyTemplate(name string) error {
	if !templateExists(name) {
		fmt.Println("Invalid language:", name)
		return nil
	}

	file, err := os.Create(dest)
	if err != nil {
		return errors.New("Unable to open destination file")
	}
	defer file.Close()

	template, err := templates.Open("templates/" + name)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(template)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}

func main() {
	languages, err := getLanguages()
	if err != nil {
		fmt.Println("Unable to read templates")
		return
	}
	args := os.Args[1:]
	if 0 == len(args) {
		fmt.Println("Choose a language to generate the template:", strings.Join(languages, ", "))
		return
	}
	err = copyTemplate(args[0])
	if err != nil {
		fmt.Println("Unable to copy template:", err)
		return
	}

	fmt.Println("Config file created at:", dest)
}
