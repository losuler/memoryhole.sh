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

func parseTemplate(templateDir string, buildDir string, filename string, c Config) {
    // Read templates.
    template, err := template.ParseFiles(templateDir + "/" + filename)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if _, err := os.Stat(buildDir); os.IsNotExist(err) {
        os.Mkdir(buildDir, 0700)
    }

    file, err := os.Create(buildDir + "/" + filename)
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

func buildTemplates(templateDir string, buildDir string, c Config) {
    files, err := ioutil.ReadDir(templateDir)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            parseTemplate(templateDir, buildDir, filename, c)
        }
    }
    fmt.Println("[INFO] Successfully generated files from templates.")
}

func main() {
    c := readConfig()

    buildTemplates("..", string(c.Footer.Version), c)
    buildTemplates("../vlsm", string(c.Footer.Version + "/vlsm"), c)

    rsync := exec.Command("rsync", "--archive",
        "../css", "../fonts", c.Footer.Version)
    if err := rsync.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    rsync = exec.Command("rsync", "--archive",
        "../vlsm/vlsm.go", "../vlsm/vlsm.service", c.Footer.Version + "/vlsm/")
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
