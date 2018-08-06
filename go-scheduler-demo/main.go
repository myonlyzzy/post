package main

import (
	"net/http"
	"log"
	"runtime"
	"fmt"
	"syscall"
)

func main() {
	fmt.Println(runtime.NumCPU())
	http.HandleFunc("/sleep", func(w http.ResponseWriter, request *http.Request) {
		//time.Sleep(time.Second * 300
		tspec := syscall.NsecToTimespec(1000 * 1000 * 1000)
		if err := syscall.Nanosleep(&tspec, &tspec); err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/echo", func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("hi"))
	})
	log.Println(http.ListenAndServe(":8070", nil))
}
