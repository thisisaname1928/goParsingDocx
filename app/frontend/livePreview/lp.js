const modeToggleButton = document.getElementById('modeToggleButton');
const body = document.body;
const filePathInput = document.getElementById('filePathInput')
const msg = document.getElementById('msg')
const questions = document.getElementById('questions')

document.addEventListener('DOMContentLoaded', () => {
    const savedMode = localStorage.getItem('theme');
    if (savedMode === 'dark') {
        body.classList.add('dark-mode');
        modeToggleButton.textContent = 'Chế độ sáng';
    } else {
        modeToggleButton.textContent = 'Chế độ tối';
    }
})

modeToggleButton.addEventListener('click', () => {
    body.classList.toggle('dark-mode');

    if (body.classList.contains('dark-mode')) {
        modeToggleButton.textContent = 'Chế độ sáng';
        localStorage.setItem('theme', 'dark');
    } else {
        modeToggleButton.textContent = 'Chế độ tối';
        localStorage.setItem('theme', 'light');
    }
});

function transTNAnswer(A, B, C, D) {
    if (A) {
        return 'A'
    }

    if (B) {
        return 'B'
    }

    if (C) {
        return 'C'
    }

    if (D) {
        return 'D'
    }
}

function transTNDSAnswer(A, B, C, D) {
    output = ""

    if (A) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (B) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (C) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    if (D) {
        output += 'Đ'
    } else {
        output += 'S'
    }

    return output
}

function transTLNAnswer(a1, a2, a3, a4) {
    return a1 + a2 + a3 + a4
}

function fetchAPI(name, obj) {
    return fetch('/LivePreview/API/' + name, { method: 'POST', body: JSON.stringify(obj) })
}

function getTnAns(sheet) {
    for (i = 0; i < 4; i++) {
        if (sheet[i]) {
            return String.fromCharCode(65 + i)
        }
    }

    return "Chua co"
}

function prepareQuestions(json) {
    if (!json.status) {
        if (json.error == "invalid docx file") {
            msg.innerText = "File không hợp lệ!"
        } else {
            msg.innerText = json.error
        }
        return;
    }
    msg.innerText = ""

    ques = ""
    for (i = 0; i < json.questions.length; i++) {
        if (json.questions[i].type == 0x12) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc loại ${json.questions[i].stype} : ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-letter">A.</div>
                            <div class="option-text">${json.questions[i].answers[0]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">B.</div>
                            <div class="option-text">${json.questions[i].answers[1]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">C.</div>
                            <div class="option-text">${json.questions[i].answers[2]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">D.</div>
                            <div class="option-text">${json.questions[i].answers[3]}</div>
                        </div>
                        <div>Đáp án: ${transTNAnswer(json.questions[i].TNAnswers[0], json.questions[i].TNAnswers[1], json.questions[i].TNAnswers[2], json.questions[i].TNAnswers[3])}</div>
                    </div>
                </div > `
        } else if (json.questions[i].type == 0x13) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc nghiệm trả lời ngắn loại ${json.questions[i].stype}: ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-text">Đáp án: ${transTLNAnswer(json.questions[i].TLNAnswers[0], json.questions[i].TLNAnswers[1], json.questions[i].TLNAnswers[2], json.questions[i].TLNAnswers[3])}</div>
                        </div>
                    </div>
                </div> `
        } else if (json.questions[i].type == 0x15) {
            ques += `<div class="question-card">
                    <div class="question-text">
                        Câu trắc nghiệm đúng sai loại ${json.questions[i].stype} : ${json.questions[i].content}
                    </div>
                    <div class="options-list">
                        <div class="option-item">
                            <div class="option-letter">A.</div>
                            <div class="option-text">${json.questions[i].answers[0]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">B.</div>
                            <div class="option-text">${json.questions[i].answers[1]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">C.</div>
                            <div class="option-text">${json.questions[i].answers[2]}</div>
                        </div>
                        <div class="option-item">
                            <div class="option-letter">D.</div>
                            <div class="option-text">${json.questions[i].answers[3]}</div>
                        </div>
                        <div>Đáp án: ${transTNDSAnswer(json.questions[i].TNAnswers[0], json.questions[i].TNAnswers[1], json.questions[i].TNAnswers[2], json.questions[i].TNAnswers[3])}</div>
                    </div>
                </div> `}
    }
    questions.innerHTML = ques;
}

async function fetchForQuestions() {
    return await fetch('/LivePreview/API/genJson', { method: 'POST', body: JSON.stringify({ path: filePathInput.value }) }).then((r) => { r.json().then((json) => prepareQuestions(json)) })
}

async function getUUID() {
    response = await fetch("/Export/API/genUUID", { method: "POST" })
    res = await response.text()

    return await res
}

async function internalUploadFileAPI(path) {
    if (!path || path.trim() === "") {
        alert("Vui lòng chọn hoặc nhập đường dẫn file trước khi xuất đề!");
        return;
    }
    UUID = await getUUID()

    res = await fetch("/LivePreview/API/internalUploadAPI", { method: "POST", body: JSON.stringify({ path: path, UUID: UUID }) })
    resJ = await res.json()

    if (!resJ.status) {
        if (resJ.msg == "ERR_FILE_NOT_FOUND") {
            alert("Đường dẫn nhập bị sai!")
        } else {
            alert("Có lỗi trong quá trình tìm file!")
        }
        return
    }

    window.location.href = window.location.origin + "/Export/Config/UUID/" + UUID + "?exportType=useDocx"
}

async function browseFile() {
    if (window.pywebview && window.pywebview.api && window.pywebview.api.select_file) {
        try {
            const selectedPath = await window.pywebview.api.select_file();
            if (selectedPath) {
                filePathInput.value = selectedPath;
                fetchForQuestions();
            }
        } catch (e) {
            console.error("Error opening OS file dialog:", e);
        }
    } else {
        alert("Tính năng chọn file qua HĐH chỉ khả dụng khi chạy ứng dụng qua Douglas Desktop (wrapper). Bạn vẫn có thể nhập đường dẫn file vào ô bên cạnh.");
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

livePreviewLoop()
async function livePreviewLoop() {
    while (true) {
        await sleep(500);
        await fetchForQuestions()
    }
}