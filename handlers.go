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
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type jsonTime time.Time

func (t jsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *jsonTime) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseFloat(string(s), 64)

	if err != nil {
		return
	}

	*(*time.Time)(t) = time.Unix(int64(q), 0)
	return
}

type PreflightRequest struct {
	SerialNo     string `json:"serial_no"`
	Hostname     string `json:"hostname"`
	SantaVersion string `json" json:"santa_version"`
	OsVersion    string `json:"os_version"`
	OsBuild      string `json:"os_build"`
	PrimaryUser  string `json:"primary_user"`
}

type Machine struct {
	MachineId    string `gorm:"column:machineid"`
	SerialNo     string `gorm:"column:serial_no"`
	Hostname     string `gorm:"column:hostname"`
	SantaVersion string `gorm:"column:santa_version"`
	OsVersion    string `gorm:"column:os_version"`
	OsBuild      string `gorm:"column:os_build"`
	PrimaryUser  string `gorm:"column:primary_user"`
}

func (c Machine) TableName() string {
	return "machines"
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
	Cursor string `json:"cursor,omitempty"`
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
	EventId         int64      `gorm:"column:eventid; primary_key:yes`
	MachineId       string     `gorm:"column:machineid"`
	FileSha1        string     `gorm:"column:sha1" json:"file_sha1"`
	FilePath        string     `gorm:"column:filepath" json:"file_path" `
	FileName        string     `gorm:"column:filename" json:"file_name"`
	ExecutionUser   string     `sql:"-"; json:execution_user`
	ExecutionTime   jsonTime   `gorm:"column:execution_time" sql:"-" json:"execution_time"`
	Decision        EventState `gorm:"column:decision" sql:"type:int;" json:"decision"`
	LoggedInUsers   []string   `sql:"-" json:"logged_in_users"`
	CurrentSessions []string   `sql:"-" json:"current_sessions"`

	FileBundleId            string `sql:"-" json:file_bundle_id`
	FileBundleName          string `sql:"-" json:file_bundle_name`
	FileBundleVersion       string `sql:"-" json:file_bundle_version`
	FileBundleVersionString string `sql:"-" json:file_bundle_version_string`

	CertificateSha1       string    `sql:"-" json:cert_sha1`
	CertificateCN         string    `sql:"-" json:cert_cn`
	CertificateOrg        string    `sql:"-" json:cert_org`
	CertificateOU         string    `sql:"-" json:cert_ou`
	CertificateValidFrom  time.Time `sql:"-" json:cert_valid_from`
	CertificateValidUntil time.Time `sql:"-" json:cert_valid_until`
}

func preFlightHandler(c *gin.Context) {
	machineid := c.Params.ByName("machineid")

	log.Printf("Preflight received from: %s\n", machineid)

	var t PreflightRequest
	c.Bind(&t)

	var machine *Machine = &Machine{
		MachineId:    machineid,
		SerialNo:     t.SerialNo,
		Hostname:     t.Hostname,
		SantaVersion: t.SantaVersion,
		OsVersion:    t.OsVersion,
		OsBuild:      t.OsBuild,
		PrimaryUser:  t.PrimaryUser,
	}

	tx := db.Begin()
	tx.Create(machine)
	tx.Commit()

	var rdr PreflightResponse
	rdr.BatchSize = 20
	rdr.UploadLogsUrl = fmt.Sprintf("%s://%s/upload/%s", c.Request.URL.Scheme, c.Request.Host, machineid)
	rdr.ClientMode = ClientModeMonitor
	c.JSON(200, rdr)
}

func postFlightHandler(c *gin.Context) {
	machineid := c.Params.ByName("machineid")

	log.Printf("Postflight received from: %s\n", machineid)

	c.JSON(200, gin.H{})
}

func ruleDownloadHandler(c *gin.Context) {
	machineid := c.Params.ByName("machineid")

	log.Printf("Sending rules to: %s\n", machineid)

	var t RuleDownloadRequest
	c.Bind(&t)

	var rdr RuleDownloadResponse
	// rdr.Rules = make([]Rule, 0)
	rdr.Rules = []*Rule{&Rule{Sha1: "d6b3853583a7dd19449275723b05b7a7e75d4529", State: RuleStateRemove, Type: RuleTypeBinary, CustomMessage: "test"}}
	c.JSON(200, rdr)
}

func eventUploadHandler(c *gin.Context) {
	machineid := c.Params.ByName("machineid")

	log.Printf("Receiving events for: %s\n", machineid)

	var t EventUploadRequest
	c.Bind(&t)

	tx := db.Begin()
	fmt.Println(tx)

	fmt.Println(t)
	fmt.Println(t.Events)

	for _, e := range t.Events {
		fmt.Println(e.ExecutionTime)
		fmt.Println(e.FilePath)
		e.MachineId = machineid
		tx.Create(e)
		// db.Exec("INSERT INTO events (sha1, filename) VALUES (?, ?)", e.FileSha1, e.FileName)
	}

	tx.Commit()

	c.JSON(500, gin.H{})

	c.JSON(200, gin.H{})
}

func uploadLogHandler(c *gin.Context) {
	machineid := c.Params.ByName("machineid")

	log.Printf("Receiving logs from: %s\n", machineid)

	mr := multipart.NewReader(c.Request.Body, "santa-sync-upload-boundary" /*params["boundary"]*/)
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

	c.JSON(200, gin.H{})
}

func viewHandler(c *gin.Context) {
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
