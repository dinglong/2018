package handler

import (
	"io"
	"net/http"

	"xregistry/frontend"
)

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp.WriteHeader(http.StatusOK)

	data, err := frontend.Asset("resource/index.html")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, err.Error())
		return
	}

	resp.Write(data)
}
