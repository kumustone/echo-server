package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var version = "V0.1"

type Config struct {
	EchoPath     string `toml:"echo_path"`
	DataPath     string `toml:"data_path"`
	DownloadPath string `toml:"download_path"`
	UseGzip      bool   `toml:"use_gzip"`
	UseChunk     bool   `toml:"use_chunk"`
}

var config Config

func main() {
	setupLogging()
	changeWorkingDirectory()

	configFile, helpFlag, versionFlag := parseFlags()

	if *helpFlag {
		printHelp()
		return
	}

	if *versionFlag {
		printVersion()
		return
	}

	loadConfig(*configFile)

	startServer()
}

func setupLogging() {
	logDir := "./log/echo-server"
	logFile := filepath.Join(logDir, "echo-server.log")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Error: Unable to create log directory: %v", err)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error: Unable to open log file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println("echo-server started")
}

func changeWorkingDirectory() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error: Unable to get executable path: %v", err)
	}
	workDir := filepath.Dir(execPath)
	if err := os.Chdir(workDir); err != nil {
		log.Fatalf("Error: Unable to change working directory: %v", err)
	}
	log.Printf("Changed working directory to %s", workDir)
}

func parseFlags() (configFile *string, helpFlag *bool, versionFlag *bool) {
	configFile = flag.String("c", "", "Path to config file")
	helpFlag = flag.Bool("h", false, "Show help information")
	versionFlag = flag.Bool("v", false, "Show version information")
	flag.Parse()
	return
}

func loadConfig(configFile string) {
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			log.Fatalf("Error: Unable to read config file: %v", err)
		}
	} else {
		defaultConfigPath := filepath.Join(".", "config.toml")
		if _, err := os.Stat(defaultConfigPath); err == nil {
			if _, err := toml.DecodeFile(defaultConfigPath, &config); err != nil {
				log.Fatalf("Error: Unable to read config file: %v", err)
			}
		} else {
			config = Config{
				EchoPath:     "echo",
				DataPath:     "data",
				DownloadPath: "download",
				UseGzip:      false,
				UseChunk:     false,
			}
		}
	}
}

func startServer() {
	server := http.Server{
		Addr: "0.0.0.0:64000",
	}

	http.HandleFunc("/help", usageHandler)
	http.HandleFunc("/data/", dataHandler)
	http.HandleFunc("/echo/", echoHandler)

	fileServer := http.FileServer(http.Dir(config.DownloadPath))
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
  -h, --help      Show help information
  -v, --version   Show version information

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
