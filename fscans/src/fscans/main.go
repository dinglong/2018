package main

import (
	"flag"
	"fscans/analyser"
	"fscans/pkg"
	"fscans/works"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var root string
	var maxRoutine int

	flag.StringVar(&root, "path", "", "The path to analyser")
	flag.IntVar(&maxRoutine, "max-routine", 0, "The max routines to analyser")
	flag.Parse()

	if len(root) == 0 {
		flag.Usage()
		return
	}

	doAnalyser(root, maxRoutine)
}

func doAnalyser(root string, maxRoutine int) {
	ws := make(chan string)
	fs := make(chan *pkg.FileHash)

	// 对目录进行遍历，遇到文件时将全路径写入works chan
	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				ws <- path
			}
			return nil
		})
		close(ws)
	}()

	worksWait := &sync.WaitGroup{}
	worksWait.Add(1)
	go func() {
		defer worksWait.Done()
		h := works.NewHashWorker(ws, fs, maxRoutine)
		h.Run()
	}()

	analyserWait := &sync.WaitGroup{}
	analyserWait.Add(1)
	go func() {
		defer analyserWait.Done()
		a := analyser.NewAnalyser(fs)
		a.Run()
		a.Dump()
	}()

	// 等待work结束后关闭fs chan，以驱动analyser结束
	worksWait.Wait()
	close(fs)

	// 等待analyser结束
	analyserWait.Wait()
}
