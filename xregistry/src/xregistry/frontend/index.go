package frontend

import "net/http"

func init() {
	http.Handle("/", http.FileServer(assetFS()))
}
