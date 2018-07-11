package works

import (
	"crypto/sha256"
	"fmt"
	"fscans/pkg"
	"io"
	"os"
	"runtime"
	"sync"
)

type hashWorker struct {
	maxRoutine int
	works      chan string // works chan关闭则worker退出
	fileHashes chan *pkg.FileHash
}

func NewHashWorker(works chan string, fileHashes chan *pkg.FileHash, maxRoutine int) *hashWorker {
	routines := maxRoutine
	if routines <= 0 {
		routines = runtime.NumCPU()
	}

	return &hashWorker{
		maxRoutine: routines,
		works:      works,
		fileHashes: fileHashes,
	}
}

func (w *hashWorker) Add(p string) {
	w.works <- p
}

// Run 启动worker，从channel中读取文件进行分析
func (w *hashWorker) Run() {
	// 用于等待所有的工作协程退出
	worksWait := sync.WaitGroup{}
	worksWait.Add(w.maxRoutine)

	for i := 0; i < w.maxRoutine; i++ {
		go func() {
			defer worksWait.Done()
			for {
				filename, more := <-w.works
				if !more {
					return
				}

				hash, err := hashFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "worker catch error %v\n", err)
					continue
				}

				w.fileHashes <- &pkg.FileHash{
					Hash:     hash,
					FileName: filename,
				}
			}
		}()
	}

	worksWait.Wait()
}

func hashFile(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte{})), nil
}
