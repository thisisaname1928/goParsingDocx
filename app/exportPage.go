package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thisisaname1928/goParsingDocx/docx"
	"github.com/thisisaname1928/goParsingDocx/dou"
)

func exportRouteRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/export/")
}

func exportRoute(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/export/index.html")
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

func exportConfigRouteRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/export/config/")
}

func exportConfigRoute(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	rawUUID := v["UUID"]
	if idx := strings.Index(rawUUID, "?"); idx != -1 {
		rawUUID = rawUUID[:idx]
	}
	cleanUUID := filepath.Base(rawUUID)
	_, e := os.Stat("./app/tests/" + cleanUUID + ".dat")

	if e != nil {
		w.Write([]byte("NO SUCH FILE OR DIR"))
		return
	}

	file, e := os.Open("./app/frontend/export/config/index.html")
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

// type ExportRequest struct {
// 	Author           string              `json:"author"`
// 	UseEncryption    bool                `json:"useEncryption"`
// 	Key              string              `json:"key"`
// 	UseTestStructure bool                `json:"useTestStructure"`
// 	TestStruct       []dou.TestStructure `json:"testStructure"`
// }

// gen a v4 uuid
func genUUID() string {
	id := uuid.New()
	return id.String()
}

type GetConfigRequest struct {
	UUID       string `json:"UUID"`
	ExportType string `json:"exportType"`
}

type ExportRequest struct {
	Status               bool                        `json:"status"`
	UUID                 string                      `json:"UUID"`
	Msg                  string                      `json:"msg"`
	TestDuration         uint64                      `json:"testDuration"`
	Author               string                      `json:"author"`
	Key                  string                      `json:"key"`
	Stype                []StypeN                    `json:"stype"`
	ExportType           string                      `json:"exportType"`
	AdditionalExportData dou.DouAdditionalExportData `json:"additionalExportData"`
}

type ExportRespone struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func exportAPI(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	// client should upload file with uuid then using ExportRequest with uuid

	switch v["NAME"] {
	case "getConfig":
		var request GetConfigRequest
		var response ExportConfigResponse
		decoder := json.NewDecoder(r.Body)
		e := decoder.Decode(&request)

		if e != nil {
			response.Status = false
			response.Msg = "invalid request"
		} else {
			switch request.ExportType {
			case "useDocx":
				response, _ = getExportConfig(request.UUID)
			case "useRawText":
				response, _ = getExportConfigRawTest(request.UUID)
			}
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(&response)
	case "export":
		var request ExportRequest
		var response ExportRespone
		decoder := json.NewDecoder(r.Body)
		e := decoder.Decode(&request)

		if e != nil {
			response.Status = false
			response.Msg = fmt.Sprintf("%v", e)
			return
		}

		var testStructure []dou.TestStructure
		// im too silly to convert this, it is the same=)))
		for _, v := range request.Stype {
			var curS dou.TestStructure
			curS.N = v.N
			curS.Stype = v.Stype
			curS.Points = v.Point
			testStructure = append(testStructure, curS)
		}

		useEncryption := false
		if request.Key != "" {
			useEncryption = true
		}

		switch request.ExportType {
		case "useDocx":
			e = dou.Export("./app/tests/"+request.UUID+".dat", "./app/tests/"+request.UUID+".dou", request.Author, request.TestDuration, true, testStructure, useEncryption, request.Key, request.AdditionalExportData)
			if e != nil {
				response.Msg = "invalid docx file"
				response.Status = false

				encoder := json.NewEncoder(w)
				encoder.Encode(&response)
				return
			}
		case "useRawText":
			b, e := os.ReadFile("./app/tests/" + request.UUID + ".dat")

			if e != nil {
				response.Msg = "invalid docx file"
				response.Status = false

				encoder := json.NewEncoder(w)
				encoder.Encode(&response)
				return
			}

			e = dou.ExportWithFluid(docx.String2Fluid(string(b)), "./app/tests/"+request.UUID+".dou", request.Author, request.TestDuration, true, testStructure, useEncryption, request.Key, request.AdditionalExportData)
			if e != nil {
				response.Msg = "invalid docx file"
				response.Status = false

				encoder := json.NewEncoder(w)
				encoder.Encode(&response)
				return
			}
		}

		response.Msg = "ok"
		response.Status = true

		encoder := json.NewEncoder(w)
		encoder.Encode(&response)
	case "saveToPath":
		var request struct {
			UUID       string `json:"UUID"`
			TargetPath string `json:"targetPath"`
		}
		var response struct {
			Status bool   `json:"status"`
			Msg    string `json:"msg"`
		}
		decoder := json.NewDecoder(r.Body)
		e := decoder.Decode(&request)
		if e != nil || request.UUID == "" || request.TargetPath == "" {
			response.Status = false
			response.Msg = "invalid parameters"
			json.NewEncoder(w).Encode(&response)
			return
		}

		cleanUUID := filepath.Base(request.UUID)
		srcPath := "./app/tests/" + cleanUUID + ".dou"
		b, e := os.ReadFile(srcPath)
		if e != nil {
			response.Status = false
			response.Msg = "file not found"
			json.NewEncoder(w).Encode(&response)
			return
		}

		e = os.WriteFile(request.TargetPath, b, 0666)
		if e != nil {
			response.Status = false
			response.Msg = fmt.Sprintf("%v", e)
			json.NewEncoder(w).Encode(&response)
			return
		}

		response.Status = true
		response.Msg = "ok"
		json.NewEncoder(w).Encode(&response)
	case "upload":
		cleanUUID := filepath.Base(r.Header.Get("uuid"))
		_ = os.MkdirAll("./app/tests", 0777)
		f, e := os.Create("./app/tests/" + cleanUUID + ".dat")
		if e != nil {
			fmt.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		dat, e := io.ReadAll(r.Body)
		if e != nil {
			fmt.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, e = f.Write(dat); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "genUUID":
		w.Write([]byte(genUUID()))
		w.Header().Add("Content-Type", "text/plain")
	}
}

type StypeN struct {
	Stype string  `json:"stype"`
	N     uint64  `json:"N"`
	Point float64 `json:"Point"`
}

type ExportConfigResponse struct {
	Status            bool     `json:"status"`
	Msg               string   `json:"msg"`
	NumberOfQuestions uint64   `json:"numberOfQuestions"`
	Stype             []StypeN `json:"stype"`
}

func search4Stype4Cfg(ques []StypeN, name string) int {
	for i, v := range ques {
		if v.Stype == name {
			return i
		}
	}

	return -1
}

// function to get some info about test
func getExportConfig(UUID string) (ExportConfigResponse, error) {
	var response ExportConfigResponse
	res, e := GenQues("./app/tests/" + UUID + ".dat")
	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprintf("%v", e)

		return response, e
	}

	response.Status = true
	response.NumberOfQuestions = uint64(len(res))
	for _, v := range res {
		curIdx := search4Stype4Cfg(response.Stype, v.Stype)
		if curIdx != -1 { // there are some questions already have that stype in response
			response.Stype[curIdx].N++
		} else {
			response.Stype = append(response.Stype, StypeN{v.Stype, 1, 0})
		}
	}

	return response, nil
}

func getExportConfigRawTest(UUID string) (ExportConfigResponse, error) {
	var response ExportConfigResponse
	resp, e := os.ReadFile("./app/tests/" + UUID + ".dat")
	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprintf("%v", e)

		return response, e
	}

	tokens := docx.Lex(docx.String2Fluid(string(resp)))
	res := docx.BetterParse(tokens)

	response.Status = true
	response.NumberOfQuestions = uint64(len(res))
	for _, v := range res {
		curIdx := search4Stype4Cfg(response.Stype, v.Stype)
		if curIdx != -1 { // there are some questions already have that stype in response
			response.Stype[curIdx].N++
		} else {
			response.Stype = append(response.Stype, StypeN{v.Stype, 1, 0})
		}
	}

	return response, nil
}

func downloadTestRoute(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	rawUUID := v["UUID"]
	if idx := strings.Index(rawUUID, "?"); idx != -1 {
		rawUUID = rawUUID[:idx]
	}
	cleanUUID := filepath.Base(rawUUID)
	path := "./app/tests/" + cleanUUID

	b, e := os.ReadFile(path)
	fmt.Println(path)
	if e != nil {
		w.Write([]byte("404 NOT FOUND!"))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+cleanUUID+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b)
}
