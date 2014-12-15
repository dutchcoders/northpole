/*
The MIT License (MIT)

Copyright (c) 2014 DutchCoders [https://github.com/dutchcoders/]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var config struct {
}

func init() {
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()

	r.HandleFunc("/postflight/{machineid}", postFlightHandler).Methods("POST")
	r.HandleFunc("/preflight/{machineid}", preFlightHandler).Methods("POST")
	r.HandleFunc("/ruledownload/{machineid}", ruleDownloadHandler).Methods("POST")
	r.HandleFunc("/eventupload/{machineid}", eventUploadHandler).Methods("POST")
	r.HandleFunc("/upload/{machineid}", uploadLogHandler).Methods("POST")
	r.HandleFunc("/", viewHandler).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	port := flag.String("port", "8080", "port number, default: 8080")
	logpath := flag.String("log", "", "")
	flag.Parse()

	if *logpath != "" {
		f, err := os.OpenFile(*logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		defer f.Close()

		log.SetOutput(f)
	}

	log.Printf("Northpole (Santa sync): server started. :\nlistening on port: %v\n", *port)
	log.Printf("---------------------------")

	// TODO: https listener
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: handlers.PanicHandler(LoveHandler(handlers.LogHandler(r, handlers.NewLogOptions(log.Printf, "_default_"))), nil),
	}

	log.Panic(s.ListenAndServe())
	log.Printf("Server stopped.")
}
