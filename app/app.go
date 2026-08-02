package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type appVersion struct {
	VersionInt      int    `json:"versionInt"`
	VersionStr      string `json:"versionStr"`
	ShouldUpdate    bool   `json:"shouldUpdate"`
	UpdatePathWin64 string `json:"updatePathWin64"`
}

func getVersion(w http.ResponseWriter, _ *http.Request) {
	var response appVersion
	encoder := json.NewEncoder(w)

	response.VersionInt = -1
	response.VersionStr = "BAD_?"

	b, e := os.ReadFile("./appVersion.json")

	if e != nil {
		encoder.Encode(response)
		return
	}

	json.Unmarshal(b, &response)

	encoder.Encode(response)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/icon.ico")
	if e != nil {
		w.Write([]byte{})
		return
	}
	defer file.Close()
	f, e := io.ReadAll(file)

	if str := detectFileExt("./app/icon.ico"); str != "" {
		w.Header().Add("Content-Type", str)
	} else {
		contentType := http.DetectContentType(f)
		w.Header().Add("Content-Type", contentType)
	}

	if e == nil {
		w.Write(f)
	}
}

func addResource(w http.ResponseWriter, r *http.Request, path string) {
	vars := mux.Vars(r)
	cleanFile := filepath.Base(vars["FILE"])
	file, e := os.Open(path + cleanFile)
	if e != nil {
		w.Write([]byte("NOT FOUND"))
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

func check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func StartApp() {
	server := mux.NewRouter()
	server.HandleFunc("/LivePreview", livePreview)
	server.HandleFunc("/LivePreview/{FILE}", livePreviewRes)
	server.HandleFunc("/LivePreview/API/{NAME}", livePreviewAPI)
	server.HandleFunc("/favicon.ico", favicon)
	server.HandleFunc("/media/{FILE}", mediaRoute)
	server.HandleFunc("/Home", homePage)
	server.HandleFunc("/Home/{FILE}", homePageRes)
	server.HandleFunc("/Export", exportRoute)
	server.HandleFunc("/Export/{FILE}", exportRouteRes)
	server.HandleFunc("/Export/API/{NAME}", exportAPI)
	server.HandleFunc("/Export/Config/{FILE}", exportConfigRouteRes)
	server.HandleFunc("/Export/Config/UUID/{UUID}", exportConfigRoute)
	server.HandleFunc("/Export/Download/UUID/{UUID}", downloadTestRoute)
	server.HandleFunc("/StartTest.TestInfo/uuid/{UUID}", testIn4)
	server.HandleFunc("/StartTest.TestInfo/{FILE}", testInfoRes)
	server.HandleFunc("/StartTest", startTest)
	server.HandleFunc("/StartTest/{FILE}", startTestRes)
	server.HandleFunc("/StartTest/API/{NAME}", startTestAPI)
	server.HandleFunc("/TutorialPage", tutorialPage)
	server.HandleFunc("/TutorialPage/{FILE}", tutorialPageRes)
	server.HandleFunc("/check", check)
	server.HandleFunc("/getVersion", getVersion)
	server.HandleFunc("/quickEditor", quickEditorRoute)
	server.HandleFunc("/quickEditor/{FILE}", quickEditorRouteRes)
	server.HandleFunc("/API/genAI", genAIAPI)
	server.HandleFunc("/API/quickPreview", quickPreviewAPI)
	server.HandleFunc("/update", updatePage)
	server.HandleFunc("/update/{FILE}", updatePageRes)
	server.HandleFunc("/downloadUpdate", downloadUpdate)
	fmt.Println("dia chi web app: http://localhost:8080/Home")
	http.ListenAndServe("localhost:8080", server)
}
