package main

import (
	"flag"
	"fmt"
	"github.com/rymis/leo-learns-prog/learnsrv/server"
	"log"
	"net/http"
	"os"
)

func main() {
	dataPath := flag.String("data", "./data", "specify directory with user data")
	wwwPath := flag.String("wwwdata", "../programming-tasks/dist", "path to static WWW data")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	srv, err := server.New(*dataPath, *wwwPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: srv,
	}
	log.Fatal(s.ListenAndServe())
}
