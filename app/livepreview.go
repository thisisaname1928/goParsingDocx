package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/docx"
)

func livePreview(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/livePreview/index.html")
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

func detectFileExt(path string) string {
	fileExtSpl := strings.Split(path, ".")
	ext := fileExtSpl[len(fileExtSpl)-1]

	switch ext {
	case "js":
		return "text/javascript"
	case "css":
		return "text/css"
	case "ico":
		return "image/x-icon"
	default:
		return ""
	}
}

func livePreviewRes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cleanFile := filepath.Base(vars["FILE"])
	file, e := os.Open("./app/frontend/livePreview/" + cleanFile)
	if e != nil {
		w.Write([]byte{})
		return
	}
	defer file.Close()
	f, e := io.ReadAll(file)

	if str := detectFileExt(cleanFile); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}

// allow upload file with path
func internalUploadAPI(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path string `json:"path"`
		UUID string `json:"UUID"`
	}

	var response struct {
		Status bool   `json:"status"`
		Msg    string `json:"msg"`
	}

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	e := decoder.Decode(&request)
	if e != nil {
		response.Msg = "ERR_FILE_NOT_FOUND"
		response.Status = false
		encoder.Encode(response)
		return
	}

	f, e := os.ReadFile(request.Path)

	if e != nil {
		response.Msg = "ERR_FILE_NOT_FOUND"
		response.Status = false
		encoder.Encode(response)
		return
	}

	_ = os.MkdirAll("./app/tests", 0777)
	e = os.WriteFile("./app/tests/"+request.UUID+".dat", f, os.FileMode(0777))
	if e != nil {
		response.Msg = "INTERNAL_ERR"
		response.Status = false
		encoder.Encode(response)
		return
	}

	response.Msg = "ok"
	response.Status = true
	encoder.Encode(response)
}

type GenJsonAPIRequest struct {
	Path string `json:"path"`
}

type GenJsonResponse struct {
	Status    bool            `json:"status"`
	Error     string          `json:"error"`
	Msg       string          `json:"msg"`
	Questions []docx.Question `json:"questions"`
}

func livePreviewAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switch vars["NAME"] {
	case "genJson":
		decoder := json.NewDecoder(r.Body)
		var request GenJsonAPIRequest
		decoder.Decode(&request)

		res, e := GenQues(request.Path)

		var response GenJsonResponse
		encoder := json.NewEncoder(w)

		if e != nil {
			response.Status = false
			response.Error = fmt.Sprintf("%v", e)
		} else {
			response.Status = true
			response.Questions = res
		}
		encoder.Encode(&response)
	case "internalUploadAPI":
		internalUploadAPI(w, r)
	}
}

func ConvertPath(path *string) {
	r := []rune(*path)
	for i := range r {
		if r[i] == '\\' {
			r[i] = '/'
		}
	}

	*path = string(r)
}

func GenQues(path string) ([]docx.Question, error) {
	ConvertPath(&path)
	// dest := strings.Split(path, "/")
	// ddest := dest[len(dest)-1]

	docx.DecompressDocxMedia(path, "./app/media/")
	fluid, e := docx.Parse2Fluid(path)
	if e != nil {
		return []docx.Question{}, e
	}

	tokens := docx.Lex(fluid)
	var ques []docx.Question = docx.BetterParse(tokens)
	return ques, nil
}

func mediaRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cleanFile := filepath.Base(vars["FILE"])
	file, e := os.Open("./app/media/" + cleanFile)
	if e != nil {
		w.Write([]byte{})
		return
	}
	defer file.Close()
	f, e := io.ReadAll(file)
	if e != nil {
		return
	}

	if str := detectFileExt("./app/media/" + cleanFile); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	w.Write(f)
}
