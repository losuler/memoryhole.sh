package main

import (
    "fmt"
    "os"
    "os/exec"
    "io/ioutil"
    "strings"
    "html/template"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Footer struct {
        Version string
        Date string
        Source string
    }
}

func readConfig() (Config) {
    var c Config

    raw, err := ioutil.ReadFile("config.yml")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    err = yaml.Unmarshal([]byte(raw), &c)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return c
}

func parseTemplate(filename string, c Config) {
    var dir string

    // Read templates.
    template, err := template.ParseFiles("../" + filename)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    dir = string(c.Footer.Version)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.Mkdir(dir, 0700)
    }

    file, err := os.Create(dir + "/" + filename)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Write output to new file.
    err = template.Execute(file, c)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    files, err := ioutil.ReadDir("..")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Read config.
    c := readConfig()

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            parseTemplate(filename, c)
        }
    }
    fmt.Println("[INFO] Successfully generated files from templates.")

    rsync := exec.Command("rsync", "--archive",
        "../css", "../fonts", c.Footer.Version)
    if err := rsync.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    gzipName := c.Footer.Version + ".tar.gz"

    gzip := exec.Command("tar", "--create",
        "--file", gzipName, c.Footer.Version)

    if err := gzip.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("[INFO] Successfully created compressed archive.")

    if err := os.RemoveAll(c.Footer.Version); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
