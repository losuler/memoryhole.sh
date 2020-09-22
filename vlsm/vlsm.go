package main

import (
    "html/template"
    "net/http"
    "os/exec"
    "fmt"
    "regexp"
    "log"
    "strings"
)

type Data struct {
    //Output string
    Output template.HTML
    Err string
}

func processInput(input string) Data {
    var data Data

    inputSlice := strings.Fields(input)

    // TODO: Check commands exist and are the correct ones (e.g. ipcalc)
    cmd := fmt.Sprintf("script --quiet -c 'ipcalc-jodies %s --s %s' | ansifilter -H -f",
    //cmd := fmt.Sprintf("script --quiet --log-out /dev/null -c 'ipcalc %s --s %s' | ansifilter -H -f",
    //cmd := fmt.Sprintf("ipcalc %s --s %s", 
                       inputSlice[0], strings.Trim(fmt.Sprint(inputSlice[1:]), "[]"))

    out, err := exec.Command("bash", "-c", cmd).Output()
    if err != nil {
        log.Fatal(err)
        data.Err = "An internal error occured."
    } else {
       //data.Output = string(out)

       out := strings.TrimSuffix(string(out), "\n")

       // Replace blue color
       out = strings.ReplaceAll(string(out), "#0000ee", "#7587a6")
       // Replace red color
       out = strings.ReplaceAll(string(out), "#cd0000", "#914343")
       // Replace pink color
       out = strings.ReplaceAll(string(out), "#cd00cd", "#91608b")
       // Replace green color
       out = strings.ReplaceAll(string(out), "#00cd00", "#5e8855")
       // Replace yellow color
       out = strings.ReplaceAll(string(out), "#cdcd00", "#a8a866")
       
       data.Output = template.HTML(out)
    }
    
    return data
}

func main() {
    var input string

    vlsm := template.Must(template.New("vlsm.html").Delims("[[", "]]").ParseFiles("vlsm.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var data Data

        if r.Method != http.MethodPost {
            vlsm.Execute(w, nil)
            return
        }

        input = r.FormValue("prompt")
        exp := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}[ \d]+$`)

        if exp.MatchString(input) {
            data = processInput(input)
        } else {
            data.Err = "Please check your input."
        }

        vlsm.Execute(w, data)
    })

    http.ListenAndServe(":6275", nil)
}
