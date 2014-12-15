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
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type PreflightRequest struct {
	SerialNo     string `json:"serial_no"`
	Hostname     string `json:hostname`
	SantaVersion string `json:"santa_version"`
	OsVersion    string `json:"os_version"`
	OsBuild      string `json:"os_build"`
	PrimaryUser  string `json:"primary_user"`
}

type PreflightResponse struct {
	BatchSize     int        `json:"batch_size"`
	UploadLogsUrl string     `json:upload_logs_url,omitempty`
	ClientMode    ClientMode `json:"client_mode"`
}

func (t *PreflightRequest) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s", t.SerialNo, t.Hostname, t.SantaVersion, t.OsVersion, t.OsBuild, t.PrimaryUser)
}

type RuleDownloadRequest struct {
	cursor string `json:"cursor"`
}

type RuleDownloadResponse struct {
	Rules  []*Rule `json:"rules"`
	Cursor string  `json:"cursor,omitempty"`
}

type EventUploadRequest struct {
	Events []Event `json:"events"`
}

type Rule struct {
	Sha1          string    `json:"sha1"`
	State         RuleState `json:"state"` // (whitelist, blacklist, silent_blacklist, remove)
	Type          RuleType  `json:"type"`  // (binary, certificate)
	CustomMessage string    `json:"custom_msg,omitempty"`
}

type Event struct {
	FileSha1        string     `json:"file_sha1"`
	FilePath        string     `json:"file_path"`
	FileName        string     `json:"file_name"`
	ExecutionUser   string     `json:execution_user`
	ExecutionTime   float64    `json:execution_time`
	Decision        EventState `json:"decision"`
	LoggedInUsers   []string   `json:"logged_in_users"`
	CurrentSessions []string   `json:"current_sessions"`

	FileBundleId            string `json:file_bundle_id`
	FileBundleName          string `json:file_bundle_name`
	FileBundleVersion       string `json:file_bundle_version`
	FileBundleVersionString string `json:file_bundle_version_string`

	CertificateSha1       string    `json:cert_sha1`
	CertificateCN         string    `json:cert_cn`
	CertificateOrg        string    `json:cert_org`
	CertificateOU         string    `json:cert_ou`
	CertificateValidFrom  time.Time `json:cert_valid_from`
	CertfiicateValidUntil time.Time `json:cert_valid_until`
}

func preFlightHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineid := vars["machineid"]

	log.Printf("Preflight received from: %s\n", machineid)

	var t PreflightRequest
	if err := ReadJSON(w, r, &t); err != nil {
		panic(err)
	}

	var rdr PreflightResponse
	rdr.BatchSize = 20
	rdr.UploadLogsUrl = fmt.Sprintf("%s://%s/upload/%s", r.URL.Scheme, r.Host, machineid)
	rdr.ClientMode = ClientModeMonitor
	if err := WriteJSON(w, r, rdr); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func postFlightHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineid := vars["machineid"]

	log.Printf("Postflight received from: %s\n", machineid)

	w.WriteHeader(http.StatusOK)
}

func ruleDownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineid := vars["machineid"]

	log.Printf("Sending rules to: %s\n", machineid)

	var t RuleDownloadRequest
	if err := ReadJSON(w, r, &t); err != nil {
		panic(err)
	}

	var rdr RuleDownloadResponse
	// rdr.Rules = make([]Rule, 0)
	rdr.Rules = []*Rule{&Rule{Sha1: "d6b3853583a7dd19449275723b05b7a7e75d4529", State: RuleStateBlacklist, Type: RuleTypeBinary}}
	if err := WriteJSON(w, r, rdr); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func eventUploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineid := vars["machineid"]

	log.Printf("Receiving events for: %s\n", machineid)
	var t EventUploadRequest
	if err := ReadJSON(w, r, &t); err != nil {
		panic(err)
	}

	for _, e := range t.Events {
		fmt.Println(e)
	}

	w.WriteHeader(http.StatusOK)
}

func uploadLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	machineid := vars["machineid"]

	log.Printf("Receiving logs from: %s\n", machineid)

	mr := multipart.NewReader(r.Body, "santa-sync-upload-boundary" /*params["boundary"]*/)
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		slurp, err := ioutil.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Part %q: %q\n", p.Header.Get("Foo"), slurp)
	}

	w.WriteHeader(http.StatusOK)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

func LoveHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-made-with", "<3 by DutchCoders")
		w.Header().Set("x-served-by", "Proudly served by DutchCoders")
		w.Header().Set("Server", "Santa Sync HTTP Server 1.0")
		h.ServeHTTP(w, r)
	}
}
