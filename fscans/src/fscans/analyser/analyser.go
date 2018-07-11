package analyser

import (
	"fmt"
	"fscans/pkg"
)

type analyser struct {
	entities   map[string][]string // hash - filename
	fileHashes chan *pkg.FileHash  // fileHashes chan 关闭则分析器退出
}

func NewAnalyser(fileHashes chan *pkg.FileHash) *analyser {
	return &analyser{
		entities:   make(map[string][]string),
		fileHashes: fileHashes,
	}
}

// Run 启动分析器，读去文件hash结果，存入map中
func (a *analyser) Run() {
	for {
		select {
		case fileHash, more := <-a.fileHashes:
			if !more {
				return
			}

			if files, exist := a.entities[fileHash.Hash]; exist {
				a.entities[fileHash.Hash] = append(files, fileHash.FileName)
			} else {
				files = []string{fileHash.FileName}
				a.entities[fileHash.Hash] = files
			}
		}
	}
}

// Dump 打印分析结果
func (a *analyser) Dump() {
	count := 0
	for h, fs := range a.entities {
		fmt.Printf("+ %s\n", h)
		for _, f := range fs {
			fmt.Printf("- %s\n", f)
			count++
		}
	}
	fmt.Printf("\nsum files: %d\n", count)
}
