package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type handler struct {
	root string
}

func (handler *handler) serveRoot(rw http.ResponseWriter, req *http.Request) {
	args := rootTemplateArgs{
		prettyPath(handler.root), "",
	}

	rootTemplate.Execute(rw, args)
}

const okMessage = `file has been uploaded<br/>`

func (handler *handler) serveOK(rw http.ResponseWriter, req *http.Request) {
	args := rootTemplateArgs{
		prettyPath(handler.root), okMessage,
	}

	rootTemplate.Execute(rw, args)
}

func (handler *handler) serveUpload(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Print("not an upload request")
		http.Error(rw, "must be an upload request", http.StatusBadRequest)
		return
	}

	file, header, err := req.FormFile("source")
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(rw, fmt.Sprintf("error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	filename := filepath.Base(header.Filename)
	dest, err := os.Create(filepath.Join(handler.root, filename))
	if err != nil {
		log.Printf("can't create destination file: %v", err)
		http.Error(rw, fmt.Sprintf("can't create destination file: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dest, file)
	if err != nil {
		log.Printf("failed to save file: %v", err)
		http.Error(rw, fmt.Sprintf("failed to save file: %v", err), http.StatusInternalServerError)
		return
	}
	dest.Close()

	http.Redirect(rw, req, "/ok", http.StatusTemporaryRedirect)
}

func prettyPath(src string) string {
	if filepath.IsAbs(src) {
		return src
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("can't get working dir: %v", err)
		return src
	}

	return filepath.Join(wd, src)
}

func main() {
	listenOn := flag.String("listen", ":9780", "listen on specified endpoint")
	root := flag.String("root", ".", "fs root to serve from")

	flag.Parse()

	handler := &handler{
		root: *root,
	}

	http.HandleFunc("/", handler.serveRoot)
	http.HandleFunc("/ok", handler.serveOK)
	http.HandleFunc("/upload", handler.serveUpload)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(*root))))

	err := http.ListenAndServe(*listenOn, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to serve: %v", err)
	}
}
