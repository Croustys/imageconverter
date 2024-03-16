package main

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sunshineplan/imgconv"
)

type ReqRes struct {
	http.ResponseWriter
	*http.Request
}

const templates = "./templates"

var formatMap = map[string]imgconv.Format{
	"PNG":  imgconv.PNG,
	"JPG":  imgconv.JPEG,
	"JPEG": imgconv.JPEG,
	"GIF":  imgconv.GIF,
	"PDF":  imgconv.PDF,
}

func handler(w http.ResponseWriter, r *http.Request) {
	rr := &ReqRes{ResponseWriter: w, Request: r}

	if r.Method == "GET" {
		rr.handleGet()
	} else if r.Method == "POST" {
		rr.handlePost()
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (r *ReqRes) handleGet() {
	http.ServeFile(r.ResponseWriter, r.Request, getPath("/index.html"))
}

func (r *ReqRes) handlePost() {
	ext := r.Request.FormValue("ext")
	file, header, err := r.Request.FormFile("image")
	if err != nil {
		http.Error(r.ResponseWriter, "Image not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	src, err := imgconv.Decode(file)
	if err != nil {
		http.Error(r.ResponseWriter, "Failed to decode image", http.StatusInternalServerError)
		return
	}

	filename := header.Filename
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))

	newFormat := getFormat(ext)
	buffer := new(bytes.Buffer)
	err = imgconv.Write(buffer, src, &imgconv.FormatOption{Format: newFormat})
	if err != nil {
		http.Error(r.ResponseWriter, "Failed to convert image", http.StatusInternalServerError)
		return
	}

	r.setFileHeaders(len(buffer.Bytes()), filenameWithoutExt+"."+strings.ToLower(ext))

	if _, err := r.ResponseWriter.Write(buffer.Bytes()); err != nil {
		http.Error(r.ResponseWriter, "Failed to write image data", http.StatusInternalServerError)
		return
	}
}
func (r *ReqRes) setFileHeaders(len int, fileName string) {
	r.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	r.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	r.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len))
}
func getPath(p string) string {
	return filepath.Join(templates, p)
}
func getFormat(p string) imgconv.Format {
	if format, ok := formatMap[strings.ToUpper(p)]; ok {
		return format
	}
	return imgconv.PNG
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Server listening on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Println("Error:", err)
	}
}
