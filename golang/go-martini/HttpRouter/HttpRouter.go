package main

import (
	"fmt"
	"runtime"
	"time"

	//"net/http"

	"github.com/shines77/martini"
	//"github.com/martini-contrib/render"
	//"github.com/odysseus/stopwatch"
)

// Run benchmark of various urls
var testUrls = []struct {
	method string
	path   string
	match  martini.RouteMatch
}{
	{"GET", "/service/candy/lollipop", martini.ExactMatch},
	{"GET", "/service/candy/gum", martini.ExactMatch},
	{"GET", "/service/candy/seg_ratta", martini.ExactMatch},
	{"GET", "/service/candy/lakrits", martini.ExactMatch},

	{"GET", "/service/shutdown", martini.ExactMatch},
	{"GET", "/", martini.ExactMatch},
	{"GET", "/some_file.html", martini.ExactMatch},
	{"GET", "/another_file.jpeg", martini.ExactMatch},
}

// Run benchmark of various urls
var testUrls2 = []struct {
	method string
	path   string
	match  martini.RouteMatch
}{
	{"GET", "/service/candy/lollipop", martini.ExactMatch},
	{"GET", "/service/candy/gum", martini.ExactMatch},
	{"GET", "/service/candy/seg_ratta", martini.ExactMatch},
	{"GET", "/service/candy/lakrits", martini.ExactMatch},

	{"GET", "/service/shutdown", martini.ExactMatch},
	{"GET", "/", martini.ExactMatch},
	{"GET", "/download/some_file.html", martini.ExactMatch},
	{"GET", "/download/another_file.jpeg", martini.ExactMatch},
}

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
	userRouted := 0
	router := martini.NewRouter()

	router.Get("/service/candy/:kind", func(params martini.Params) {
		userRouted++
	})

	router.Get("/service/shutdown", func() {
		userRouted++
	})

	router.Get("/", func() {
		userRouted++
	})

	router.Get("/download/:filename", func(params martini.Params) {
		userRouted++
	})

	for i, r := range router.GetAllRoutes() {
		fmt.Printf("[%3d]: %-20s %-8s %-30s\n", i,
			r.GetName(), r.Method(), r.Pattern())
	}
	fmt.Printf("\n")

	for _, urls := range testUrls {
		matched := false
		for _, route := range router.GetAllRoutes() {
			match, params := route.Match(urls.method, urls.path)
			if match != martini.NoMatch {
				fmt.Printf("Matched: (%v, %-30v) - (%v), (%v)\n",
					urls.method, urls.path, match, params)
				matched = true
				break
			}
		}
		if !matched {
			fmt.Printf("Not matched: (%v, %-30v)\n", urls.method, urls.path)
		}
		runtime.Gosched()
	}
	fmt.Printf("\n")
}

func Benchmark_Routing() {
	userRouted := 0
	router := martini.NewRouter()

	router.Get("/service/candy/:kind", func(params martini.Params) {
		userRouted++
	})

	router.Get("/service/shutdown", func() {
		userRouted++
	})

	router.Get("/", func() {
		userRouted++
	})

	router.Get("/:filename", func(params martini.Params) {
		userRouted++
	})

	kMaxIterators := 1000000
	fmt.Printf("kMaxIterators = %d\n\n", kMaxIterators)

	startTotalTime := time.Now()
	for _, urls := range testUrls {
		matched := false
		startTime := time.Now()
		for i := 0; i < kMaxIterators; i++ {
			for _, route := range router.GetAllRoutes() {
				match, _ := route.Match(urls.method, urls.path)
				if match != martini.NoMatch {
					matched = true
					break
				}
			}
			if matched {
				userRouted++
			}
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

func Benchmark_Routing2() {
	userRouted := 0
	router := martini.NewRouter()

	router.Get("/service/candy/:kind", func(params martini.Params) {
		userRouted++
	})

	router.Get("/service/shutdown", func() {
		userRouted++
	})

	router.Get("/", func() {
		userRouted++
	})

	router.Get("/download/:filename", func(params martini.Params) {
		userRouted++
	})

	kMaxIterators := 1000000
	fmt.Printf("kMaxIterators = %d\n\n", kMaxIterators)

	startTotalTime := time.Now()
	for _, urls := range testUrls2 {
		matched := false
		startTime := time.Now()
		for i := 0; i < kMaxIterators; i++ {
			for _, route := range router.GetAllRoutes() {
				match, _ := route.Match(urls.method, urls.path)
				if match != martini.NoMatch {
					matched = true
					break
				}
			}
			if matched {
				userRouted++
			}
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
	fmt.Printf("Martini 1.0 - HttpRouter Test\n")
	fmt.Printf("\n")

	UnitTest_Routing()
	Benchmark_Routing()
	Benchmark_Routing2()
}
