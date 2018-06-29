package analyser

import (
	"fmt"
	"os"

	"fscans/pkg"
	"sync"
)

type analyser struct {
	entities   map[string][]string // hash - filename
	fileHashes chan *pkg.FileHash
	errs       chan error
}

func NewAnalyser(fileHashes chan *pkg.FileHash, errs chan error) *analyser {
	return &analyser{
		entities:   make(map[string][]string),
		fileHashes: fileHashes,
		errs:       errs,
	}
}

func (a *analyser) Run(wait *sync.WaitGroup) {
	defer wait.Done()

	for {
		select {
		case fileHash, ok := <-a.fileHashes:
			if !ok {
				return
			}

			if files, exist := a.entities[fileHash.Hash]; exist {
				a.entities[fileHash.Hash] = append(files, fileHash.FileName)
			} else {
				files = []string{fileHash.FileName}
				a.entities[fileHash.Hash] = files
			}

		case err := <-a.errs:
			fmt.Fprintf(os.Stderr, "analyser catch error %v\n", err)
		}
	}
}

func (a *analyser) Dump() {
	count := 0
	for h, fs := range a.entities {
		fmt.Printf("%s\n", h)
		for _, f := range fs {
			fmt.Printf("\t %s\n", f)
			count++
		}
	}
	fmt.Printf("sum files: %d\n", count)
}
