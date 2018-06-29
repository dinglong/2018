package main

import (
	"os"
	"path/filepath"

	"fscans/analyser"
	"fscans/pkg"
	"fscans/works"
	"sync"
)

func main() {
	root := "E:\\temp"

	ws := make(chan string)
	fs := make(chan *pkg.FileHash)
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	a := analyser.NewAnalyser(fs, errs)
	go a.Run(wg)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				ws <- path
			}
			return nil
		})
		close(ws)
	}()

	h := works.NewHashWork(0, ws, fs, errs)
	h.Run(wg)

	close(fs)
	wg.Wait()

	a.Dump()
}
