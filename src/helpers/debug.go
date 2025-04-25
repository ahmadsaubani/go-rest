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

// Dd - Dump and exit. This is similar to Laravel's dd helper.
// It dumps the given values and exits the program.
func Dd(data ...interface{}) {
	for _, d := range data {
		spew.Dump(d)
	}
	os.Exit(1)
}

// Dump is similar to Dd, but it does not exit the program. It
// simply dumps the given values to the console.
func Dump(data ...interface{}) {
	for _, d := range data {
		spew.Dump(d)
	}
}

// DumpJSON is similar to Dump, but it formats the output as JSON.
func DumpJSON(data ...interface{}) {
	for _, d := range data {
		fmt.Printf("%+v\n", d)
	}
}

// DdLog logs the given data to a file with the current date as the filename
// and also prints it to the console. It creates the log directory if it does
// not exist. After logging, the function exits the program.

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

// DumpLog logs the given data to a file with the current date as the filename
// and also prints it to the console. It creates the log directory if it does
// not exist. After logging, the function does not exit the program.
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
