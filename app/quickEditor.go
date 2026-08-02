package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/thisisaname1928/goParsingDocx/docx"
	genaiwrapper "github.com/thisisaname1928/goParsingDocx/genAIWrapper"
)

func quickEditorRouteRes(w http.ResponseWriter, r *http.Request) {
	addResource(w, r, "./app/frontend/quickEditor/")
}

func quickEditorRoute(w http.ResponseWriter, r *http.Request) {
	file, e := os.Open("./app/frontend/quickEditor/index.html")
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

type GeminiPart struct {
	Texts []string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiCandinates struct {
	Content GeminiContent `json:"content"`
}

type GeminiResponse struct {
	Candinates GeminiCandinates `json:"candinates"`
}

type Root struct {
	Candidates    []Candidate   `json:"candidates"`
	UsageMetadata UsageMetadata `json:"usageMetadata"`
	ModelVersion  string        `json:"modelVersion"`
	ResponseID    string        `json:"responseId"`
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
	AvgLogprobs  float64 `json:"avgLogprobs"`
}

type Content struct {
	Parts []Part `json:"parts"`
	Role  string `json:"role"`
}

type Part struct {
	Text string `json:"text"`
}

type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

func Call4AI(prompt string, key string) (string, error) {
	reqBodyObj := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
	}
	jsonBytes, err := json.Marshal(reqBodyObj)
	if err != nil {
		return "", err
	}

	requestBody := bytes.NewBuffer(jsonBytes)
	resquest, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent", requestBody)
	if err != nil {
		return "", err
	}

	resquest.Header.Add("X-goog-api-key", key)
	resquest.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(resquest)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error status: %s", res.Status)
	}

	decoder := json.NewDecoder(res.Body)
	var response Root
	if err := decoder.Decode(&response); err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty AI response candidates")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func genAIAPI(w http.ResponseWriter, r *http.Request) {
	var request struct {
		NumberOfQuesTN   int    `json:"numberOfQuesTN"`
		NumberOfQuesTNDS int    `json:"numberOfQuesTNDS"`
		NumberOfQuesTLN  int    `json:"numberOfQuesTLN"`
		Content          string `json:"content"`
	}

	var response struct {
		Status bool   `json:"status"`
		Msg    string `json:"msg"`
		Res    string `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	e := decoder.Decode(&request)

	if e != nil {
		response.Status = false
		response.Msg = "ERR_BAD_REQUEST"
		encoder.Encode(response)
	}

	response.Res, e = genaiwrapper.AutoGenContent("Tôi là giáo viên tôi muốn tạo " + fmt.Sprint(request.NumberOfQuesTN) + " câu trắc nghiệm, " + fmt.Sprint(request.NumberOfQuesTNDS) + " câu trắc nghiệm đúng sai, " + fmt.Sprint(request.NumberOfQuesTLN) + " câu trả lời ngắn (câu trả lời tối đa 4 kí tự và là 1 số) cho học sinh bao gồm  nội dung trong kiến thức về " + request.Content + " hãy trả về cho tôi theo đúng template ko giải thích hay nói gì thêm làm lệch template. Template có dạng Tất cả câu hỏi bắt đầu bằng chữ 'Câu' và phải bắt buộc cho 'C' viết hoa.Các đáp án của câu hỏi trắc nghiệm bắt đầu bằng A. B. C. D. và bắt buộc phải có 4 đáp án.Các đáp án của câu hỏi trắc nghiệm đúng sai bắt đầu bằng a) b) c) d) và bắt buộc phải có 4 đáp án.Đáp án của trắc nghiệm trả lời ngắn nằm ở phía sau cặp từ tln_ans: ...Về đáp án đúng của câu trắc nghiệm và trắc nghiệm đúng sai được đặt dấu * liền trước chữ cái đại diện cho đáp án, Ví dụ: *A., *a),....Nội dung câu hỏi có hỗ trợ định dạng bằng html, và bắt buộc đoạn code định dạng không có bất kì định dạng nào từ trình soạn thảo văn bản.Loại câu hỏi được đặt trong ngoặc vuông phía sau số thứ tự câu hỏi. VD: Câu 1 [LOAI_1]: là câu hỏi có loại là LOAI_1. Lưu ý giữa những câu hỏi không nên có dòng trống, Các câu hỏi liền nhau. Khi chèn các dấu <, > mà không dùng cho mục đích định dạng html nên sử dụng &lt và &gt.")

	if e != nil {
		response.Status = false
		response.Msg = fmt.Sprint(e)
	} else {
		response.Status = true
		response.Msg = "ok"
	}
	encoder.Encode(response)
}

func quickPreviewAPI(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Content string `json:"content"`
	}

	var response GenJsonResponse

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	e := decoder.Decode(&request)

	if e != nil {
		response.Status = false
		response.Msg = "ERR_BAD_REQUEST"
		encoder.Encode(response)
		return
	}

	fluids := docx.String2Fluid(request.Content)

	tokens := docx.Lex(fluids)
	ques := docx.BetterParse(tokens)

	response.Questions = ques
	response.Msg = "ok"
	response.Status = true
	encoder.Encode(response)
}
