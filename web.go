package main

import (
    "net/http"
    "net/url"
    "os"
    "log"
    "fmt"
    "bufio"
    "regexp"
    "strconv"
)

const timeFormat = "20060102"
var getDT = regexp.MustCompile(`^(DT(?:START|END):)(\d{4})(.*)?$`)

func timewarpServer(w http.ResponseWriter, req *http.Request) {
    values, err := url.ParseQuery(req.URL.RawQuery)

    if err != nil {
        log.Panic("Couldn't decode querystring %v", err)
    }

    ical := values.Get("ical")
    if ical == "" {
        w.WriteHeader(400)
        fmt.Fprintln(w, "missing ?ical=")
        log.Println("Couldn't decode querystring %v", err)
        return
    }

    resp, err := http.Get(ical)
    defer resp.Body.Close()
    if err != nil {
        log.Printf("http error %v\n", err)
        w.WriteHeader(502)
        fmt.Fprintln(w, "couldn't retrieve your URL")
        return
    }

    w.Header().Add("Content-Type", "text/calendar")
    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()
        match := getDT.FindStringSubmatch(line)
        if match != nil {
            i, err := strconv.Atoi(match[2])
            if err == nil {
                // One year in future
                fmt.Fprintf(w, "%s%d%s\r\n", match[1], i + 1, match[3])
                continue
            }
        }
        fmt.Fprintln(w, line)
    }
}

func main() {
    http.HandleFunc("/", timewarpServer)

    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}
