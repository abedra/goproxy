package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	cache  *LRU
	log    *os.File
}

func (proxy *Proxy) handle(writer http.ResponseWriter, request *http.Request) {
	request.Header.Set("X-GoProxy", "Success")

	if _, ok := proxy.cache.Get(request.RemoteAddr); ok {
		log.Printf("Cache HIT for %s", request.RemoteAddr)
	} else {
		log.Printf("Cache MISS for %s", request.RemoteAddr)
		proxy.cache.Add(request.RemoteAddr, true)
	}

	proxy.proxy.ServeHTTP(writer, request)
}

func ValidateTarget(target string) (*url.URL, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, errors.New("Target is not a valid URL")
	}
	return url, nil
}

func ValidatePort(port int) (string, error) {
	if port > 0 && port < 65537 {
		return fmt.Sprintf(":%d", port), nil
	}
	return "", errors.New("Port is invalid")
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	targetPtr := flag.String("target", "http://localhost:8080", "Proxy Target")
	portPtr := flag.Int("port", 8888, "Proxy Port")
	logPathPtr := flag.String("logdir", "logs", "Logfile Directory")
	cacheSizePtr := flag.Int("cachesize", 0, "Sets the size of the cache")
	flag.Parse()

	if _, err := os.Stat(*logPathPtr); os.IsNotExist(err) {
		err := os.Mkdir(*logPathPtr, 0755)
		if err != nil {
			fmt.Println("Error creating log directory:", err)
			os.Exit(1)
		}
	}

	logFileName := fmt.Sprintf("%s/goproxy.log", *logPathPtr)
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	url, err := ValidateTarget(*targetPtr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	port, err := ValidatePort(*portPtr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	proxy := &Proxy{
		target: url,
		proxy:  httputil.NewSingleHostReverseProxy(url),
		cache:  NewCache(*cacheSizePtr),
		log:    logFile,
	}

	http.HandleFunc("/", proxy.handle)

	log.Printf("Starting goproxy on %s. Target URL is %s\n", port, url)

	http.ListenAndServe(port, Log(http.DefaultServeMux))
}
