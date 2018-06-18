package handler

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
)

func LoadFileHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := req.ParseForm(); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, err.Error())
		return
	}

	fn := req.FormValue("file")
	log.Printf("file： %s\n", fn)

	fd, err := os.Open(fn)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, err.Error())
		return
	}
	defer fd.Close()

	resp.WriteHeader(http.StatusOK)
	io.Copy(resp, bufio.NewReader(fd))
}
