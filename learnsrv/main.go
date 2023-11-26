package main

import (
    "github.com/rymis/leo-learns-prog/learnsrv/server"
    "fmt"
    "os"
    "net/http"
    "log"
)

func main() {
    srv, err := server.New("./data", "../wwwroot")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    s := &http.Server{
        Addr:           ":8080",
        Handler:        srv,
    }
    log.Fatal(s.ListenAndServe())
}
