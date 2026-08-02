const testContent = document.getElementById("testContent")
const summitBtn = document.getElementById("summitTest")
const testResult = document.getElementById("testResult")
let isDone = false
let testModeInited = false

async function checkIfTestDone() {
    res = await fetch("/api/getTestStatus", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
    jsonRes = await res.json()

    return jsonRes.status
}

function showElement(id) { document.getElementById(id).classList.remove("hidden-element") }

function hideElement(id) {
    document.getElementById(id).classList.add("hidden-element")
}

function warnPopup() {
    showElement("popup-ask")
}

document.addEventListener('DOMContentLoaded', async () => {
    a = this.window.location.href.split("/")
    uuid = a[a.length - 1]
    isDone = await checkIfTestDone()

    if (!isDone) {
        showElement('popup-ask')
    } else {
        startTest()
    }

    document.getElementById('okBut').addEventListener('click', async () => {
        if (!testModeInited) startTest()

        if (await checkIfTestDone()) {
            location.reload()
        }
    })
})

async function startTest() {
    isDone = await checkIfTestDone()

    hideElement('popup-ask')
    await getTest(uuid)
    if (!isDone) {
        showElement('timer')
        showElement('summitTest')
        setUpTestMode(uuid)
        setUpTimer(glbTestsvr.test.startTime, glbTestsvr.test.testDuration)


        document.getElementById("schoolName").innerHTML = `Trường: ${await getSchoolName()}`
        document.getElementById("candinateName").innerHTML = `Học sinh: ${glbTestsvr.test.name}`
        document.getElementById("testName").innerHTML = `Bài thi: ${await getTestName()}`
    } else {
        showElement('timer')
        showElement('testResult')
        res = await fetch("/api/getPoint", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
        jsonRes = await res.json()

        // document.getElementById("test").innerHTML = JSON.stringify(jsonRes)

        bgTime = new Date(glbTestsvr.test.startTime)
        edTime = new Date(glbTestsvr.test.endTime)
        duration = edTime - bgTime

        testResult.innerHTML = `<h1>Điểm: ${jsonRes.point}</h1><p>Số câu đúng hoàn toàn: ${jsonRes.trueQuesCount}/${questions.length}<br><br>Tên học sinh: ${glbTestsvr.test.name}<br>Lớp: ${glbTestsvr.test.class}
        <br><br>
        Thời gian làm bài: ${Math.trunc(duration / 60000)} phut ${Math.round(((duration % 60000) / 1000))} giay
        <br>Thời gian bắt đầu: ${bgTime.toLocaleTimeString("en-us", {
            hour: '2-digit',
            minute: '2-digit', hour12: false
        })}<br>Thời gian kết thúc: ${edTime.toLocaleTimeString("en-us", {
            hour: '2-digit',
            minute: '2-digit', hour12: false
        })}<br><br>Mã bài làm: ${uuid}</p>`

        loadUpTrueAns()

        setEndTime()
    }



    // load question sheet data
    loadUpAnsSheet()
}

async function isAdmin() {
    res = await fetch("/api/getTest", { method: "GET" })
    textRes = res.text()

    if (textRes == 'true') {
        return true
    } else {
        return false
    }
}

async function okAdminNoReload() {
    if (await isAdmin()) {

    } else {
        location.reload()
    }
}

function loadUpTrueAns() {
    for (i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) {
            const e = getTNQuesElement(i)
            a = questions[i].TNAnswers
            e.innerHTML += "<br>Đáp án đúng: " + transTNAnswer(a)
        }
        else if (questions[i].type == 0x15) {
            const e = getTNDSQuesElement(i)
            a = questions[i].TNAnswers
            e.innerHTML += "<br>Đáp án đúng: " + transTNDSAnswer(a)
        } else if (questions[i].type == 0x13) {
            const e = getTLNQuesElement(i)
            a = questions[i].TLNAnswers
            e.innerHTML += "<br>Đáp án đúng: " + transTLNAnswers(a)
        }
    }
}

async function chooseTNOption(i, ans, shouldUpdate) {
    // check if update client side or both client and server
    if (shouldUpdate) {
        result = await updateAnswerSheet(i, [String.fromCharCode('A'.charCodeAt(0) + ans), '', '', ''])
        //result = await updateAnswerSheet2(i, ans, String.fromCharCode('A'.charCodeAt(0) + ans))

        if (!result) {
            return
        }
    }

    if (ans > 3) {
        return;
    }
    const item = document.getElementById(`QUES.${i}.TN.${ans}`)
    if (item.classList.contains("option-item-highlighted")) { // no need to reset
        return;
    }


    item.classList.replace("option-item", "option-item-highlighted")

    for (j = 0; j < 4; j++) {
        if (j != ans) {
            document.getElementById(`QUES.${i}.TN.${j}`).classList.replace("option-item-highlighted", "option-item")
        }
    }
}

async function chooseTNDSOption(i, ansI, value, shouldUpdate) {
    rightAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.R`)
    wrongAns = document.getElementById(`QUES.${i}.TNDS.${ansI}.W`)

    if (value) {

        if (shouldUpdate) {
            result = await updateAnswerSheet2(i, ansI, "T")

            if (!result) {
                return
            }
        }

        rightAns.classList.replace("tnds-ans-r", "tnds-ans-r-highlighted")
        wrongAns.classList.replace("tnds-ans-w-highlighted", "tnds-ans-w")
    } else {
        if (shouldUpdate) {
            result = await updateAnswerSheet2(i, ansI, "F")

            if (!result) {
                return
            }
        }

        wrongAns.classList.replace("tnds-ans-w", "tnds-ans-w-highlighted")
        rightAns.classList.replace("tnds-ans-r-highlighted", "tnds-ans-r")
    }
}

async function chooseTLNAnswer(index, answerIndex) {
    inp = document.getElementById(`QUES.${index}.TLN.${answerIndex}`)

    data = inp.value
    // add space to it
    if (data == "") {
        data = " "
    }

    result = await updateAnswerSheet2(index, answerIndex, data)
    if (!result) {
        // reset value
        inp.value = ""
    }
}

function getTNQuesElement(i) {
    return document.getElementById(`QUES.${i}.TN`)
}

function getTLNQuesElement(i) {
    return document.getElementById(`QUES.${i}.TLN`)
}

function getTNDSQuesElement(i) {
    return document.getElementById(`QUES.${i}.TNDS`)
}

let questions
let glbTestsvr

async function renderTest(testsvr) {
    questions = testsvr.test.questions
    glbTestsvr = testsvr

    part1 = document.getElementById('part1')
    part2 = document.getElementById('part2')
    part3 = document.getElementById('part3')

    part1ShouldShow = false
    part2ShouldShow = false
    part3ShouldShow = false

    part1Idx = 0
    part2Idx = 0
    part3Idx = 0

    // render questions
    for (i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) { // TN
            part1ShouldShow = true
            part1Idx++
            part1.innerHTML += `
<div class="question-card" id="QUES.${i}.TN">
    <div class="question-text">
        Câu ${part1Idx} (Trắc nghiệm): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TN.0" onclick="chooseTNOption(${i}, 0, true)">
            <div class="option-letter">A.</div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.1" onclick="chooseTNOption(${i}, 1, true)">
            <div class="option-letter">B.</div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.2" onclick="chooseTNOption(${i}, 2, true)">
            <div class="option-letter">C.</div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TN.3" onclick="chooseTNOption(${i}, 3, true)">
            <div class="option-letter">D.</div>
            <div class="option-text">${questions[i].answers[3]}</div>
        </div>
    </div>
</div>`
        } else if (questions[i].type == 0x13) {
            part3ShouldShow = true
            part3Idx++
            part3.innerHTML += `    
    <div class="question-card" id="QUES.${i}.TLN">
        <div class="question-text">
            Câu ${part3Idx} (Trắc nghiệm trả lời ngắn): ${questions[i].content}
        </div>
        <div class="TLN-input-container">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 0)" id="QUES.${i}.TLN.0">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 1)" id="QUES.${i}.TLN.1">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 2)" id="QUES.${i}.TLN.2">
            <input class="square-input" type="text" maxlength="1" oninput="chooseTLNAnswer(${i}, 3)" id="QUES.${i}.TLN.3">
        </div>
    </div>`
        } else if (questions[i].type == 0x15) {
            part2ShouldShow = true
            part2Idx++
            part2.innerHTML += `
<div class="question-card" id="QUES.${i}.TNDS">
    <div class="question-text" >
        Câu ${part2Idx} (Trắc nghiệm đúng sai): ${questions[i].content}
    </div>
    <div class="options-list">
        <div class="option-item" id="QUES.${i}.TNDS.0">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.0.R" onclick="chooseTNDSOption(${i}, 0, true, true)" >D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.0.W" onclick="chooseTNDSOption(${i}, 0, false, true)">S</button>
            <div class="option-letter">a) </div>
            <div class="option-text">${questions[i].answers[0]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.1">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.1.R" onclick="chooseTNDSOption(${i}, 1, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.1.W" onclick="chooseTNDSOption(${i}, 1, false, true)">S</button>
            <div class="option-letter">b) </div>
            <div class="option-text">${questions[i].answers[1]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.2">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.2.R" onclick="chooseTNDSOption(${i}, 2, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.2.W" onclick="chooseTNDSOption(${i}, 2, false, true)">S</button>
            <div class="option-letter">c) </div>
            <div class="option-text">${questions[i].answers[2]}</div>
        </div>
        <div class="option-item" id="QUES.${i}.TNDS.3">
            <button class="tnds-ans-r" id="QUES.${i}.TNDS.3.R" onclick="chooseTNDSOption(${i}, 3, true, true)">D</button>
            <button class="tnds-ans-w" id="QUES.${i}.TNDS.3.W" onclick="chooseTNDSOption(${i}, 3, false, true)">S</button>
            <div class="option-letter">d) </div>
            <div class="option-text">${questions[i].answers[3]}</div>
        </div>
    </div>
</div>`
        }
    }

    if (!part1ShouldShow)
        hideElement('part1')

    if (!part2ShouldShow)
        hideElement('part2')

    if (!part3ShouldShow)
        hideElement('part3')
}

async function getTest(uuid) {
    res = await fetch("/api/getTest", { method: "POST", body: JSON.stringify({ uuid: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return
    }

    await renderTest(jsonRes)

    summitBtn.addEventListener("click", async function () {
        const ok = await doneTest()
        if (ok) {
            location.reload()
        }
    })
}

// load reload answer sheet
async function loadUpAnsSheet() {
    res = await fetch("/api/getCurrentAnsSheet", { method: "POST", body: JSON.stringify({ UUID: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        alert(jsonRes.msg)
        return false
    }

    for (i = 0; i < questions.length; i++) {
        // TN

        if (questions[i].type == 18) {
            if (jsonRes.ansSheet[i][0] != "") {
                chooseTNOption(i, jsonRes.ansSheet[i][0].charCodeAt(0) - 'A'[0].charCodeAt(0), false)
            }
        } else if (questions[i].type == 21) {
            for (j = 0; j < 4; j++) {
                if (jsonRes.ansSheet[i][j] != "") {
                    chooseTNDSOption(i, j, jsonRes.ansSheet[i][j] == "T", false)
                }
            }
        } else if (questions[i].type == 19) {
            for (j = 0; j < 4; j++) {
                if (jsonRes.ansSheet[i][j] != "") {
                    // manually copy
                    inp = document.getElementById(`QUES.${i}.TLN.${j}`)
                    inp.value = jsonRes.ansSheet[i][j]
                }
            }
        }
    }
}

async function doneTest() {
    res = await fetch("/api/handleDoneTest", { method: "POST", body: JSON.stringify({ UUID: uuid }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        if (jsonRes.msg == "TEST_ACCESS_DENIED") {
            alert("Bạn đã bị tước quyền làm bài!")
        } else {
            alert(`Lỗi nội bộ: ${jsonRes.msg}`)
        }
        return false
    }

    return true
}

// update single answer in answerSheet array
async function updateAnswerSheet2(index, answerIndex, data) {
    res = await fetch("/api/updateAnswer", { method: "POST", body: JSON.stringify({ UUID: uuid, index: index, answerIndex: answerIndex, data: data, shouldClear: false }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        if (jsonRes.msg == "OUT_OF_TIME") {
            okAdminNoReload()
            return
        }
        if (jsonRes.msg != "TEST_SESSION_LOCKED") { alert(jsonRes.msg) }
        return false
    }

    return true
}

async function updateAnswerSheet(i, answers) {
    // add data field just to patch
    res = await fetch("/api/updateAnswer", { method: "POST", body: JSON.stringify({ UUID: uuid, index: i, answerSheet: answers, data: "?", shouldClear: true }) })
    jsonRes = await res.json()

    if (!jsonRes.status) {
        if (jsonRes.msg == "OUT_OF_TIME") {
            okAdminNoReload()
            return
        }
        if (jsonRes.msg != "TEST_SESSION_LOCKED") { alert(jsonRes.msg) }
        return false
    }

    return true
}

function transTNAnswer(a) {
    for (j = 0; j < 4; j++) {
        if (a[j]) {
            return String.fromCharCode('A'.charCodeAt(0) + j)
        }
    }
}

function transTNDSAnswer(a) {
    res = ''
    for (j = 0; j < 4; j++) {
        if (a[j]) {
            res += 'D'
        } else
            res += 'S'
    }

    return res
}

function transTLNAnswers(a) {
    res = ''
    for (j = 0; j < 4; j++) {
        res += a[j]
    }

    return res
}

function getTNAnswer(th) {
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TN.${i}`)
        if (item.classList.contains("option-item-highlighted")) {
            if (i == 0) {
                return 'A'
            }

            if (i == 1) {
                return 'B'
            }

            if (i == 2) {
                return 'C'
            }

            if (i == 3) {
                return 'D'
            }
        }
    }

    return ''
}

function getTNDSAnswer(th) {
    // prevent smth bad
    ans = ['', '', '', '']
    for (i = 0; i < 4; i++) {
        item = document.getElementById(`QUES.${th}.TNDS.${i}.R`)
        if (item.classList.contains("tnds-ans-r-highlighted")) {
            ans[i] = 'r'
        }

        item = document.getElementById(`QUES.${th}.TNDS.${i}.W`)
        if (item.classList.contains("tnds-ans-w-highlighted")) {
            ans[i] = 'w'
        }
    }

    return ans
}

function getTLNAnswer(th) {
    ans = ['', '', '', '']

    ans[0] = document.getElementById(`QUES.${th}.TLN.0`).value
    ans[1] = document.getElementById(`QUES.${th}.TLN.1`).value
    ans[2] = document.getElementById(`QUES.${th}.TLN.2`).value
    ans[3] = document.getElementById(`QUES.${th}.TLN.3`).value

    // remove space
    for (i = 0; i < 4; i++) {
        if (ans[0] == ' ') {
            ans[0] = ''
        }
    }

    return ans
}

function getAnswer() {
    result = []
    for (let i = 0; i < questions.length; i++) {
        if (questions[i].type == 0x12) {
            result.push(getTNAnswer(i))
        } else if (questions[i].type == 0x15) {
            result.push(getTNDSAnswer(i))
        } else if (questions[i].type == 0x13) {
            result.push(getTLNAnswer(i))
        }
    }

    return result
}

async function getTestName() {
    res = await fetch("/api/getTestName", { method: "GET" })
    return await res.text()
}

async function getSchoolName() {
    res = await fetch("/api/getSchoolName", { method: "GET" })
    return await res.text()
}