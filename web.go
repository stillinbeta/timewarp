package main

import (
    "net/http"
    "net/url"
    "os"
    "log"
    "fmt"
)

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

    fmt.Fprintf(w, "You submitted %v\n", ical)
}

func main() {
    http.HandleFunc("/", timewarpServer)

    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}
