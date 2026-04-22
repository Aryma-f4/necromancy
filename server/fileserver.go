package server

import (
	"log"
	"net/http"
	"path/filepath"
)

type FileServer struct {
	Addr string
	Dir  string
}

func NewFileServer(addr, dir string) *FileServer {
	return &FileServer{
		Addr: addr,
		Dir:  dir,
	}
}

func (fs *FileServer) Start() {
	absPath, err := filepath.Abs(fs.Dir)
	if err != nil {
		log.Fatalf("Invalid directory: %v", err)
	}

	fileHandler := http.FileServer(http.Dir(absPath))
	mux := http.NewServeMux()
	mux.Handle("/", fileHandler)

	log.Printf("Starting HTTP File Server on %s, serving %s\n", fs.Addr, absPath)
	if err := http.ListenAndServe(fs.Addr, mux); err != nil {
		log.Fatalf("File server failed: %v", err)
	}
}
