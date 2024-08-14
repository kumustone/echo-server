package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//func ReadResponse(r *bufio.Reader, req *Request) (*Response, error) {
//	tp := textproto.NewReader(r)
//	resp := &Response{
//		Request: req,
//	}
//
//	// Parse the first line of the response.
//	line, err := tp.ReadLine()
//	if err != nil {
//		if err == io.EOF {
//			err = io.ErrUnexpectedEOF
//		}
//		return nil, err
//	}
//	proto, status, ok := strings.Cut(line, " ")
//	if !ok {
//		return nil, badStringError("malformed HTTP response", line)
//	}
//	resp.Proto = proto
//	resp.Status = strings.TrimLeft(status, " ")
//
//	statusCode, _, _ := strings.Cut(resp.Status, " ")
//	if len(statusCode) != 3 {
//		return nil, badStringError("malformed HTTP status code", statusCode)
//	}
//	resp.StatusCode, err = strconv.Atoi(statusCode)
//	if err != nil || resp.StatusCode < 0 {
//		return nil, badStringError("malformed HTTP status code", statusCode)
//	}
//	if resp.ProtoMajor, resp.ProtoMinor, ok = ParseHTTPVersion(resp.Proto); !ok {
//		return nil, badStringError("malformed HTTP version", resp.Proto)
//	}
//
//	// Parse the response headers.
//	mimeHeader, err := tp.ReadMIMEHeader()
//	if err != nil {
//		if err == io.EOF {
//			err = io.ErrUnexpectedEOF
//		}
//		return nil, err
//	}
//	resp.Header = Header(mimeHeader)
//
//	fixPragmaCacheControl(resp.Header)
//
//	err = readTransfer(resp, r)
//	if err != nil {
//		return nil, err
//	}
//
//	return resp, nil
//}

func readLines(data []byte) [][]byte {
	var result [][]byte
	var lineStart int

	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lineEnd := i
			if i > 0 && data[i-1] == '\r' {
				lineEnd = i - 1
			}
			result = append(result, data[lineStart:lineEnd])
			lineStart = i + 1
		}
	}

	// 处理最后一行没有换行符的情况
	if lineStart < len(data) {
		lineEnd := len(data)
		if data[lineEnd-1] == '\r' {
			lineEnd--
		}
		result = append(result, data[lineStart:lineEnd])
	}

	return result
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Join(config.EchoPath, strings.TrimPrefix(r.URL.Path, "/echo/"))

	log.Printf("Echo handler reading file: %s", filename)
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	x1 := bytes.SplitN(fileContent, []byte("\n\n"), 2)
	x2 := bytes.SplitN(fileContent, []byte("\r\n\r\n"), 2)

	var x [][]byte
	if len(x2) >= len(x1) {
		x = x2
	} else {
		x = x1
	}

	lines := readLines(x[0])

	// 优化异常处理
	if len(lines) == 0 {
		http.Error(w, "Malformed HTTP response", http.StatusInternalServerError)
		return
	}

	// 解析状态行并设置到 w 中
	statusLine := string(lines[0])
	statusParts := strings.SplitN(statusLine, " ", 3)
	if len(statusParts) < 3 {
		http.Error(w, "Malformed status line", http.StatusInternalServerError)
		return
	}

	statusCode, err := strconv.Atoi(statusParts[1])
	if err != nil {
		http.Error(w, "Invalid status code", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)

	// 处理 headers 并写入到 w 中，跳过 Content-Length
	for _, line := range lines[1:] {
		headerParts := bytes.SplitN(line, []byte(": "), 2)
		if len(headerParts) == 2 && !bytes.EqualFold(headerParts[0], []byte("Content-Length")) {
			w.Header().Set(string(headerParts[0]), string(headerParts[1]))
		}
	}

	// 写入 response body
	if len(x) > 1 && len(x[1]) > 0 {
		w.Write(x[1])
	}
}
