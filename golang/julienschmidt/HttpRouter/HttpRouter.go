package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	//"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

// Run benchmark of various urls
var testUrls = []struct {
	method string
	path   string
}{
	{"GET", "/service/candy/lollipop"},
	{"GET", "/service/candy/gum"},
	{"GET", "/service/candy/seg_ratta"},
	{"GET", "/service/candy/lakrits"},

	{"GET", "/service/shutdown"},
	{"GET", "/"},
	{"GET", "/download/some_file.html"},
	{"GET", "/download/another_file.jpeg"},
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func getMilliSeconds(startTime time.Time, stopTime time.Time) float64 {
	// Returns the elapsed stopwatch time in milliseconds
	// time.Duration
	return float64(stopTime.Sub(startTime).Nanoseconds()) / 1000000
}

func getSeconds(startTime time.Time, stopTime time.Time) float64 {
	// Returns the elapsed stopwatch time in milliseconds
	// time.Duration
	return float64(stopTime.Sub(startTime).Nanoseconds()) / 1000000000
}

func UnitTest_Routing() {
	routed := false
	var req_params httprouter.Params
	req_params = nil

	router := httprouter.New()

	router.GET("/service/candy/:kind",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			routed = true
			req_params = params
		})

	router.GET("/service/shutdown",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			routed = true
			req_params = params
		})

	router.GET("/",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			routed = true
			req_params = params
		})

	router.GET("/download/:filename",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			routed = true
			req_params = params
		})

	w := new(mockResponseWriter)

	for _, urls := range testUrls {
		routed = false
		req_params = nil
		req, _ := http.NewRequest(urls.method, urls.path, nil)
		router.ServeHTTP(w, req)
		if routed {
			fmt.Printf("Matched: (%v, %-30v) - (%v)\n", urls.method, urls.path, req_params)
		} else {
			fmt.Printf("Dismatch: (%v, %-30v)\n", urls.method, urls.path)
		}
		runtime.Gosched()
	}
	fmt.Printf("\n")
}

func Benchmark_Routing() {
	userRouted := 0

	router := httprouter.New()

	router.GET("/service/candy/:kind",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			userRouted++
		})

	router.GET("/service/shutdown",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			userRouted++
		})

	router.GET("/",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			userRouted++
		})

	router.GET("/download/:filename",
		func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			userRouted++
		})

	w := new(mockResponseWriter)

	kMaxIterators := 1000000
	fmt.Printf("kMaxIterators = %d\n\n", kMaxIterators)

	startTotalTime := time.Now()
	for _, urls := range testUrls {
		startTime := time.Now()
		for i := 0; i < kMaxIterators; i++ {
			req, _ := http.NewRequest(urls.method, urls.path, nil)
			router.ServeHTTP(w, req)
		}
		stopTime := time.Now()
		elapsedTime := getMilliSeconds(startTime, stopTime)

		fmt.Printf("Matched: (%s, %-30s) - (%10d), (%9.3f ms)\n",
			urls.method, urls.path, userRouted, elapsedTime)
		runtime.Gosched()
	}
	fmt.Printf("\n")

	stopTotalTime := time.Now()
	totalElapsedTime := getSeconds(startTotalTime, stopTotalTime)

	fmt.Printf("userRouted = %d\n\n", userRouted)
	fmt.Printf("Total elapsed time: %7.3f second(s)\n", totalElapsedTime)
	fmt.Printf("\n")
}

func main() {
	runtime.GOMAXPROCS(1)

	fmt.Printf("\n")
	fmt.Printf("/github.com/julienschmidt/httprouter - HttpRouter Test\n")
	fmt.Printf("\n")

	UnitTest_Routing()
	Benchmark_Routing()
}
