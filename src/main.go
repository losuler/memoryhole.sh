package main

import (
    //"fmt"
    "os"
    "io/ioutil"
    "strings"
    "html/template"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Footer struct {
        Version string
        Date string
    }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readConfig() (Config) {
    var c Config
    
    raw, err := ioutil.ReadFile("src/config.yml")
    check(err)

    err = yaml.Unmarshal([]byte(raw), &c)
    check(err)

    return c
}

func parseTemplate(filename string) {
	var dir string

	// Read templates.
    template, err := template.ParseFiles(filename)
    check(err)

	dir = "out"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	file, err := os.Create(dir + "/" + filename)
	check(err)

	// Read config.
    c := readConfig()

	// Write output to new file.
    err = template.Execute(file, c)
    check(err)
}

func main() {
    files, err := ioutil.ReadDir(".")
    check(err)

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            parseTemplate(filename)
        }
    }
}
