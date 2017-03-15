package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

type Download struct {
	Url  string
	Type string
}

var downloadChannel chan Download

func isAcceptableType(dtype string) bool {
	acceptableTypes := [3]string{"AUDIO", "VIDEO", "ALL"}
	for _, t := range acceptableTypes {
		if t == dtype {
			return true
		}
	}

	return false
}

func queueDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("'%s' not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "'url' param is empty or not specified", http.StatusBadRequest)
		return
	}

	dtype := r.FormValue("type")
	if dtype == "" {
		http.Error(w, "'type' param is empty or not specified", http.StatusBadRequest)
		return
	}

	if !isAcceptableType(dtype) {
		http.Error(w, fmt.Sprintf("'%s' not an acceptable download type", dtype), http.StatusBadRequest)
		return
	}

	download := Download{
		Url:  url,
		Type: dtype,
	}

	go func() {
		downloadChannel <- download
	}()

	report("Queued", download.Url, download.Type)
	fmt.Fprintf(w, "Download queued")
}

func ydownload(url string, audio bool) {
	args := []string{"-c"}
	if audio {
		args = append(args, "-x")
	}
	args = append(args, url)
	cmd := exec.Command("youtube-dl", args...)
	cmd.Dir = "/tmp/"
	out, err := cmd.Output()
	var dtype string
	if audio {
		dtype = "AUDIO"
	} else {
		dtype = "VIDEO"
	}
	var status string
	if err != nil {
		status = "Failed"
	} else {
		status = "Succeeded"
	}
	report(status, url, dtype)
	if err != nil {
		fmt.Printf("Error: (%s) %s\n", err, out)
	}
}

func downloadLoop() {
	for {
		download := <-downloadChannel
		report("Started", download.Url, download.Type)
		if download.Type == "VIDEO" {
			ydownload(download.Url, false)
		} else if download.Type == "AUDIO" {
			ydownload(download.Url, true)
		} else {
			ydownload(download.Url, false)
			ydownload(download.Url, true)
		}
	}
}

func report(status string, url string, dtype string) {
	fmt.Printf("Download %s: %s (%s)\n", status, url, dtype)
}

func main() {
	downloadChannel = make(chan Download)
	go downloadLoop()
	http.HandleFunc("/", queueDownload)
	http.ListenAndServe(":8080", nil)
}
