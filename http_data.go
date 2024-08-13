package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(config.DataPath, strings.TrimPrefix(r.URL.Path, "/data/"))
	log.Printf("data handler reading file: %s", filename)

	switch filepath.Ext(filename) {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript;charset=utf8")
	case ".json":
		w.Header().Set("Content-Type", "application/json;charset=utf8")
	case ".html":
		w.Header().Set("Content-Type", "text/html;charset=utf8")
	default:
		w.Header().Set("Content-Type", "text/plain;charset=utf8")
	}

	if config.UseChunk && config.UseGzip {
		sendGzippedChunkedResponse(w, filename)
	} else if config.UseChunk {
		sendChunkedResponse(w, filename)
	} else if config.UseGzip {
		sendGzippedResponse(w, filename)
	} else {
		http.ServeFile(w, r, filename)
	}
}

func sendChunkedResponse(w http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
	}
}

func sendGzippedResponse(w http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	defer gz.Close()

	_, err = io.Copy(gz, file)
	if err != nil {
		http.Error(w, "Error compressing file", http.StatusInternalServerError)
		return
	}
}

func sendGzippedChunkedResponse(w http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	defer gz.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			gz.Write(buf[:n])
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading or compressing file", http.StatusInternalServerError)
			return
		}
	}
}
