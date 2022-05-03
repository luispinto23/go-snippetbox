package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Application struct to hold the dependencies for the web application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Define a new command line flag with the name 'addr' and a default value of ':4000'. Also add a simple description of what the flag does.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command line flags.
	flag.Parse()

	// Create a new logger to write infotmation to the console (STDOUT).
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a new logger to write error messages to the console (STDERR).
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create a new application struct and pass in the two loggers.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Initialize a new http.Server struct setting the Addr and Handler fields
	// to the values we just defined and the errorLog field to the our custom errorLogger
	srv := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
