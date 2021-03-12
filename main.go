package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// POS - POS
const POS string = "."

// DIR - DIR
const DIR string = "/"

// PATH - PATH
const PATH string = "files"

// HTTP - HTTP
const HTTP string = "http://ansible.example.com"

func getFileName(ext string) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	datetime := time.Now().Format("200612")

	dstPath := PATH + DIR + datetime
	dstFile := dstPath + DIR + datetime + strconv.Itoa(random.Int()) + POS + strings.Replace(ext, DIR, POS, 1)

	_, err := os.Stat(dstPath)
	res := os.IsNotExist(err)
	if res == true {
		os.MkdirAll(dstPath, os.ModePerm)
	}

	return dstFile
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	fileName := getFileName(handler.Header.Get("Content-Type"))

	fp, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer fp.Close()

	size, err := io.Copy(fp, file)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	go callUploadResult(w, fileName, size)
}

func callUploadResult(w http.ResponseWriter, fileName string, size int64) {
	var list = make(map[string]string)
	list["image"] = HTTP + DIR + fileName
	list["size"] = strconv.FormatInt(size/1024, 10) + "KB"
	list["action"] = "call-upload-result"
	jsonStr, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
}

func main() {
	http.HandleFunc("/", uploadFile)
	http.ListenAndServe("0.0.0.0:8090", nil)
}
