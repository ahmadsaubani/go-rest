package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "  ",
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		ContinueOnMethod:        true,
	}
}

// Dd (Dump and Die)
func Dd(data ...interface{}) {
	for _, d := range data {
		spew.Dump(d)
	}
	os.Exit(1)
}

// Dump without exiting
func Dump(data ...interface{}) {
	for _, d := range data {
		spew.Dump(d)
	}
}

// Dump as JSON-like
func DumpJSON(data ...interface{}) {
	for _, d := range data {
		fmt.Printf("%+v\n", d)
	}
}

// DdLog - Log to file and exit
func DdLog(data ...interface{}) {
	logData := ""
	for _, d := range data {
		s := spew.Sdump(d)
		fmt.Println(s)
		logData += s + "\n"
	}

	// Create log folder if not exists
	if _, err := os.Stat("src/storage/logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	// Get today's date for file name
	filename := fmt.Sprintf("src/storage/logs/%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer file.Close()

	logger := log.New(file, "[DdLog] ", log.LstdFlags)
	logger.Println(logData)

	os.Exit(1)
}

// DumpLog - Log to file and print to console without exiting

// helpers.DumpLog("User trying to login:", user)
func DumpLog(data ...interface{}) {
	logData := ""
	for _, d := range data {
		s := spew.Sdump(d)
		fmt.Println(s)
		logData += s + "\n"
	}

	// Pastikan folder logs ada
	if _, err := os.Stat("src/storage/logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	// Nama file log berdasarkan tanggal
	filename := fmt.Sprintf("src/storage/logs/%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	logger := log.New(file, "[DumpLog] ", log.LstdFlags)
	logger.Println(logData)
}
