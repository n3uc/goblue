package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"goblue/pkg/bluefile"
)

var (
	basePath = "/api/v1"
	lsPath   = basePath + "/dir"
	fnPath   = basePath + "/file/"

	Version   = "development"
	GitCommit = "development"
	BuildTime = "development"
)

type JsonError struct {
	Err string `json:"error"`
}

type JsonFiles struct {
	Files []string `json:"files"`
}

// set cors header
func setCORS(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == "OPTIONS" {
			return
		}
		f(w, r)
	}
}

// add a path to the api,  this setups any common params like CORS
func addPath(mux *http.ServeMux, p string, f func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(p, setCORS(f))
}

func headers(w http.ResponseWriter, r *http.Request) {
	// Get requested file name
	fn := filepath.Base(r.URL.Path)
	bf, err := bluefile.New(fn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonError{Err: fmt.Sprintf("%v", err)})
		log.Printf("could not load header for %v, %v\n", fn, err)
		return
	}

	fmt.Fprintf(w, "%v", bf)
}

func lsFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(".")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonError{Err: fmt.Sprintf("unable to list current directory: %v", err)})
		log.Printf("unable to list current directory: %v\n", err)
		return
	}

	fileSlice := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			fileSlice = append(fileSlice, file.Name())
		}
	}

	json.NewEncoder(w).Encode(JsonFiles{Files: fileSlice})
}

func version(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "BuildVersion: %s\nBuild Time: %s\nGitCommit: %s\n", Version, BuildTime, GitCommit)
	fmt.Fprintf(w, "\nValid if running in Kubernetes\n")
	fmt.Fprintf(w, "  VERSION:   "+os.Getenv("VERSION")+"\n")
	fmt.Fprintf(w, "  POD_NAME:  "+os.Getenv("POD_NAME")+"\n")
	fmt.Fprintf(w, "  NODE_NAME: "+os.Getenv("NODE_NAME")+"\n")
	fmt.Fprintf(w, "  IMAGE:     "+os.Getenv("IMAGE")+"\n")
}

func startApi(port int) {
	mux := http.NewServeMux()
	addPath(mux, "/version", version)
	addPath(mux, fnPath, headers)
	addPath(mux, lsPath, lsFiles)

	add := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(add, mux))
}

func main() {
	log.Printf("goblue BuildVersion: %v  Build Time: %v  GitCommit: %v\n", Version, BuildTime, GitCommit)

	port := flag.Int("p", 9580, "the port number use for hosting the server")
	dir := flag.String("d", ".", "the directory to use as root for reference to bluefiles")
	flag.Parse()

	log.Printf("port: %v   dir: %v\n", *port, *dir)

	err := os.Chdir(*dir)
	if err != nil {
		log.Fatalf("unable to change directory to %v, %v", *dir, err)
	}

	startApi(*port)
}
