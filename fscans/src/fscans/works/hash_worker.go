package works

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"fscans/pkg"
)

type hashWork struct {
	maxRoutine int
	works      chan string
	fileHashes chan *pkg.FileHash
	errs       chan error
}

func NewHashWork(maxRoutine int, works chan string, fileHashes chan *pkg.FileHash, errs chan error) *hashWork {
	routines := maxRoutine
	if routines <= 0 {
		routines = runtime.NumCPU()
	}

	return &hashWork{
		maxRoutine: routines,
		works:      works,
		fileHashes: fileHashes,
		errs:       errs,
	}
}

func (w *hashWork) Add(p string) {
	w.works <- p
}

func (w *hashWork) Run(wait *sync.WaitGroup) {
	defer wait.Done()

	wg := sync.WaitGroup{}
	wg.Add(w.maxRoutine)

	for i := 0; i < w.maxRoutine; i++ {
		go func() {
			defer wg.Done()
			for {
				filename, ok := <-w.works
				if !ok {
					return
				}

				hash, err := hashFile(filename)
				if err != nil {
					w.errs <- err
					continue
				}

				w.fileHashes <- &pkg.FileHash{
					Hash:     hash,
					FileName: filename,
				}
			}
		}()
	}

	wg.Wait()
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
