package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(config.EchoPath, strings.TrimPrefix(r.URL.Path, "/echo/"))

	log.Printf("Echo handler reading file: %s", filename)
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(bytes.NewReader(fileContent))

	//在某些场景下，比如windows这里的'\n'可不可能是 \r\n
	statusLine, err := reader.ReadString('\n')
	if err != nil || len(statusLine) < 2 {
		http.Error(w, "Invalid HTTP response", http.StatusInternalServerError)
		return
	}

	statusCode, err := strconv.Atoi(strings.Fields(statusLine)[1])
	if err != nil {
		http.Error(w, "Invalid status code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	newlineCount := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading headers", http.StatusInternalServerError)
			return
		}

		// 检查是否是空行
		if line == "\r\n" {
			newlineCount++
			if newlineCount == 2 {
				break
			}
		} else {
			newlineCount = 0
			headerParts := strings.SplitN(line, ": ", 2)
			if len(headerParts) == 2 {
				w.Header().Set(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
			}
		}
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "Error writing response body", http.StatusInternalServerError)
		return
	}
}
