package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type timerValue struct {
	HMS          string `json:"hms"`
	HMSIndicator string `json:"hms_indicator"`
	Seconds      int    `json:"seconds"`
	Over         bool   `json:"over,omitempty"`
	Type         int    `json:"type"`
	Running      bool   `json:"running"`
}

var dm bool

func runWebServer(bindAddr string, bindPort int, darkmode bool, ctx context.Context) {
	dm = darkmode
	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))
	http.HandleFunc("/timer/", timerValueHandler)
	http.HandleFunc("/", webHandler)
	if bindAddr == "" {
		// binding to all addresses
		ips, err := findMyIPs()
		if err != nil {
			log.Fatalf("Unable to get my IP addresses: %s", err)
		}
		for _, ip := range ips {
			log.Printf("Webserver listening on %s:%d", ip, bindPort)
		}
	} else {
		log.Printf("Webserver listening on %s:%d", bindAddr, bindPort)
	}

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddr, bindPort), nil))
	}()
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{filepath.Join("templates", "base.html")}
	if dm {
		files = append(files, filepath.Join("templates", "darkmode.html"))
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, TimersAsSlice(Timers))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func timerValueHandler(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query()
	//log.Println(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	out := json.NewEncoder(w)
	path := strings.SplitN(r.URL.Path[1:], "/", -1)
	//log.Println(path)

	if len(path) != 2 {
		log.Println("invalid path in timerValueHandler")
		http.Error(w, "invalid parameters; name not specified", http.StatusBadRequest)
		return
	}
	t, ok := Timers[path[1]]
	if !ok {
		log.Printf("timer name '%s' not found", path[1])
		http.Error(w, "timer name not found", http.StatusNotFound)
		return
	}
	var tv timerValue
	tv.HMS = t.HMS()
	tv.HMSIndicator = t.HMSIndicator()
	tv.Seconds = t.Seconds()
	tv.Over = t.Over()
	tv.Type = t.Type()
	tv.Running = t.Running()
	err := out.Encode(tv)
	if err != nil {
		log.Fatalf("Unable to encode response: %s", err)
		http.Error(w, "Unable to encode response", http.StatusInternalServerError)
	}
}
