package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	version = "V0.1"
)

func main() {
	// 定义命令行参数
	helpFlag := flag.Bool("h", false, "显示帮助信息")
	versionFlag := flag.Bool("v", false, "显示版本信息")
	helpFlagLong := flag.Bool("help", false, "显示帮助信息")
	versionFlagLong := flag.Bool("version", false, "显示版本信息")

	flag.Parse()

	// 处理命令行参数
	if *helpFlag || *helpFlagLong {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag || *versionFlagLong {
		printVersion()
		os.Exit(0)
	}

	server := http.Server{
		Addr: "0.0.0.0:64000",
	}

	// 路径处理
	http.HandleFunc("/help", usageHandler)
	http.HandleFunc("/data/", dataHandler)
	http.HandleFunc("/echo/", echoHandler)

	fileServer := http.FileServer(http.Dir("./download"))
	http.Handle("/download/", http.StripPrefix("/download/", fileServer))

	fmt.Println("Server is running at http://0.0.0.0:64000")
	log.Fatal(server.ListenAndServe())
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("no IP address found")
}

func printHelp() {
	ip, err := getLocalIP()
	if err != nil {
		fmt.Println("Error: Unable to get local IP address.")
		os.Exit(1)
	}
	fmt.Printf(`Usage:
  -h, --help      显示帮助信息
  -v, --version   显示版本信息

Examples:
  curl -v  http://%s:64000/data/1.json
  curl -v  http://%s:64000/echo/login
  curl -v  http://%s:64000/data/1.json?chunk=true

`, ip, ip, ip)
}

func printVersion() {
	fmt.Println("Version:", version)
}

func usageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Usage: [options]\n\nOptions:\n  -h, --help\t\tShow this help message\n  -v, --version\t\tShow version information")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	filename := "./data/" + strings.TrimPrefix(r.URL.Path, "/data/")
	log.Printf("data handler reading file: %s", filename)

	if strings.HasPrefix(filename, "./data/") {
		fmt.Println("####1")
		ext := filepath.Ext(filename)
		switch ext {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript;charset=utf8")
		case ".json":
			w.Header().Set("Content-Type", "application/json;charset=utf8")
		case ".html":
			w.Header().Set("Content-Type", "text/html;charset=utf8")
		default:
			w.Header().Set("Content-Type", "text/plain;charset=utf8")
		}

		//TODO: 同时设置chunk和gzip的情况下，那么就返回gzip的chunk传输方式，而不是2选1
		if r.URL.Query().Get("chunk") == "true" {
			sendChunkedResponse(w, filename)
		} else if r.URL.Query().Get("gzip") == "true" {
			sendGzippedResponse(w, filename)
		} else {
			http.ServeFile(w, r, filename)
		}
	} else {
		fmt.Println("####2")
		//TODO: 返回not found file : {filename}
		http.NotFound(w, r)
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

func echoHandler(w http.ResponseWriter, r *http.Request) {
	//filename := "." + strings.TrimPrefix(r.URL.Path, "/echo/")
	//fileContent, err := os.ReadFile(filename)

	//todo: 打印当前路径
	filename := "./echo/" + strings.TrimPrefix(r.URL.Path, "/echo/")
	currentDir, _ := os.Getwd()
	log.Printf("Current directory: %s", currentDir)
	log.Printf("Echo handler reading file: %s", filename)
	fileContent, err := os.ReadFile(filename)

	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	re := regexp.MustCompile(`\r?\n\r?\n`)
	parts := re.Split(string(fileContent), 2)
	if len(parts) != 2 {
		http.Error(w, "File format is incorrect", http.StatusInternalServerError)
		return
	}
	headersPart, bodyPart := parts[0], parts[1]

	headersLines := strings.Split(headersPart, "\n")
	for _, line := range headersLines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "HTTP/") {
			statusCodeStr := strings.Fields(line)[1]
			statusCode, err := strconv.Atoi(statusCodeStr)
			if err != nil {
				http.Error(w, "Invalid status code", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(statusCode)
		} else {
			headerParts := strings.SplitN(line, ": ", 2)
			//TODO : 这里报错http: superfluous response.WriteHeader call from main.echoHandler (main.go:222)
			if len(headerParts) != 2 {
				http.Error(w, "Invalid header format", http.StatusInternalServerError)
				return
			}
			w.Header().Set(headerParts[0], headerParts[1])
		}
	}

	if r.URL.Query().Get("chunk") == "true" {
		sendChunkedBody(w, bodyPart)
	} else if r.URL.Query().Get("gzip") == "true" {
		sendGzippedBody(w, bodyPart)
	} else {
		w.Write([]byte(bodyPart))
	}
}

func sendChunkedBody(w http.ResponseWriter, body string) {
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	buf := bytes.NewBufferString(body)
	chunkSize := 1024
	for {
		chunk := make([]byte, chunkSize)
		n, err := buf.Read(chunk)
		if n > 0 {
			w.Write(chunk[:n])
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
	}
}

func sendGzippedBody(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	defer gz.Close()

	_, err := gz.Write([]byte(body))
	if err != nil {
		http.Error(w, "Error compressing body", http.StatusInternalServerError)
		return
	}
}
