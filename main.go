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
	"crypto/tls"
	_ "flag"
	_ "fmt"

	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "log"
	"net/http"
	"os"
)

var db gorm.DB

var config struct {
	DSN string
}

func init() {
	config.DSN = os.Getenv("DSN")

	var err error
	db, err = gorm.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	/*
		rand.Seed(time.Now().UTC().UnixNano())

		r := mux.NewRouter()

		r.HandleFunc("/postflight/{machineid}", postFlightHandler).Methods("POST")
		r.HandleFunc("/preflight/{machineid}", preFlightHandler).Methods("POST")
		r.HandleFunc("/ruledownload/{machineid}", ruleDownloadHandler).Methods("POST")
		r.HandleFunc("/eventupload/{machineid}", eventUploadHandler).Methods("POST")
		r.HandleFunc("/upload/{machineid}", uploadLogHandler).Methods("POST")
		r.HandleFunc("/", viewHandler).Methods("GET")

		r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

		// http.Handle("/", handlers.PanicHandler(LoveHandler(handlers.LogHandler(r, handlers.NewLogOptions(log.Printf, "_default_"))), nil))
		http.Handle("/", r)
	*/
}

func ClientCertificate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate client certificate
		// log.Println(c.Request.TLS.PeerCertificates[0])
	}
}

func main() {
	port := flag.String("port", "8080", "port number, default: 8080")
	sslport := flag.String("ssl-port", "8443", "port number, default: 8443")
	cert_file := flag.String("cert-file", "cert.pem", "")
	key_file := flag.String("key-file", "key.pem", "")
	// logpath := flag.String("log", "", "")

	flag.Parse()

	router := gin.Default()
	router.Use(ClientCertificate())

	api := router.Group("")
	api.POST("/postflight/:machineid", postFlightHandler)
	api.POST("/preflight/:machineid", preFlightHandler)
	api.POST("/ruledownload/:machineid", ruleDownloadHandler)
	api.POST("/eventupload/:machineid", eventUploadHandler)
	api.POST("/upload/:machineid", uploadLogHandler)
	api.POST("/", viewHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: router,
	}

	tlsserver := &http.Server{
		Addr:    fmt.Sprintf(":%s", *sslport),
		Handler: router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
	}

	go func() {
		// run ssl goroutine
		if err := tlsserver.ListenAndServeTLS(*cert_file, *key_file); err != nil {
			panic(err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	/*


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
	*/
}
