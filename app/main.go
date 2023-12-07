package main

import (
	"flag"
	"fmt"
	"github.com/rymis/leo-learns-prog/learnsrv/server"
	"github.com/webview/webview_go"
	"net"
	"net/http"
	"os"
)

func main() {
	dataPath := flag.String("data", "./data", "specify directory with user data")
	wwwPath := flag.String("wwwdata", "../programming-tasks/dist", "path to static WWW data")
	debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	srv, err := server.New(*dataPath, *wwwPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s := &http.Server{
		Handler: srv,
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("http://localhost:%d", port)

	go func() {
		s.Serve(listener)
	}()

	app := webview.New(*debug)
	defer app.Destroy()

	app.SetTitle("Leo learns programming")
	app.Navigate(addr)
	app.SetSize(800, 600, webview.HintNone)
	app.Run()
}
