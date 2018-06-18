package action

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func SplitCol(infile, outfile string, unit int) {
	in, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	count := 0
	contents := make([][]string, 0)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		strings.Trim(string(line), " \r\n\t")

		index := count % unit
		if len(contents) == index {
			contents = append(contents, make([]string, 0))
		}

		contents[index] = append(contents[index], string(line))
		count++
	}

	dumpContents(contents, writer)
	log.Printf("game over\n")
}

func dumpContents(contents [][]string, w io.Writer) {
	for _, row := range contents {
		line := strings.Join(row, ",")
		_, err := fmt.Fprint(w, line+"\r\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
