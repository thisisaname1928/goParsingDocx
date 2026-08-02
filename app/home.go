package app

import (
	"io"
	"net/http"
	"os"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/home/index.html")
	if e != nil {
		w.Write([]byte{})
		return
	}
	defer file.Close()
	f, e := io.ReadAll(file)

	if e == nil {
		w.Write(f)
	}
}

func homePageRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/home/")
}
