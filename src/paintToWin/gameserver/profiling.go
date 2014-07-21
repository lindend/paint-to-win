package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

func initProfiling(logfile string) error {
	if logfile != "" {
		if f, err := os.Create(logfile); err != nil {
			return err
		} else {
			return pprof.StartCPUProfile(f)
		}
	} else {
		go func() {
			http.ListenAndServe("localhost:8086", nil)
		}()
		return nil
	}
}

func stopProfiling() {
	pprof.StopCPUProfile()
}
