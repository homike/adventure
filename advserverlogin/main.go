package main

import (
	"cuttleserver/loginserver/login"
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func main() {
	// HTTP Login
	httpServer := &http.Server{}
	mux := regHandler()
	httpServer = &http.Server{
		Addr:         "0.0.0.0:9100",
		Handler:      mux,
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 30,
	}

	httpServer.ListenAndServe()
}

func regHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/Login/Api/Fishluv", login.FishluvLogin)

	return mux
}
