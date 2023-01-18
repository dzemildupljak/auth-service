package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DebugLogger exported
var DebugLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func init() {
	absPath, err := filepath.Abs("./log")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}

	debugLog, err := os.OpenFile(absPath+"/debug-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	errorLog, err := os.OpenFile(absPath+"/error-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	DebugLogger = log.New(debugLog, "Debug Logger:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	ErrorLogger = log.New(errorLog, "Error Logger:\t", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

const requestIDKey = "REQUEST_ID"

func ReqLoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				fmt.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
				DebugLogger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()

			next.ServeHTTP(w, r)
		})
	}
}
