package server

import (
	"log"
	"net/http"
	"path/filepath"
)

type FileServer struct {
	Addr   string
	Dir    string
	Prefix string
}

func NewFileServer(addr, dir, prefix string) *FileServer {
	return &FileServer{
		Addr:   addr,
		Dir:    dir,
		Prefix: prefix,
	}
}

func (fs *FileServer) Start() {
	absPath, err := filepath.Abs(fs.Dir)
	if err != nil {
		log.Fatalf("Invalid directory: %v", err)
	}

	fileHandler := http.FileServer(http.Dir(absPath))
	mux := http.NewServeMux()
	prefix := fs.Prefix
	if prefix == "" {
		prefix = "/"
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix != "/" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	if prefix == "/" {
		mux.Handle("/", fileHandler)
	} else {
		mux.Handle(prefix, http.StripPrefix(prefix, fileHandler))
	}

	log.Printf("Starting HTTP File Server on %s%s, serving %s\n", fs.Addr, prefix, absPath)
	if err := http.ListenAndServe(fs.Addr, mux); err != nil {
		log.Fatalf("File server failed: %v", err)
	}
}
