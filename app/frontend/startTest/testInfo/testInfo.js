let uuid = "NONE"
const justAButton = document.getElementById("justAButton")
const exportCsv = document.getElementById("exportCsv")
let curIP = ""

document.addEventListener('DOMContentLoaded', () => {
    uuid = window.location.pathname.replace("/StartTest.TestInfo/uuid/", "")

    updateTestInfo()

    document.getElementById("showChart").addEventListener('click', () => {
        showElement("charPopup")
        createChart(document.getElementById("chartBro"))
    })
})

async function viewLoop() {
    while (true) {
        updateTestInfo()
        await delay(1000)
    }
}

exportCsv.addEventListener('click', async () => {
    csvDat = await exportFileCSV()

    const fileContent = csvDat.csvData;
    const filename = "diem_thi.csv";
    const mimeType = "text/csv";

    const blob = new Blob([fileContent], { type: mimeType });

    const downloadLink = document.createElement("a")
    downloadLink.download = filename
    downloadLink.href = window.URL.createObjectURL(blob);

    downloadLink.click()
    window.URL.revokeObjectURL(downloadLink.href);
})


async function getTestInfo() {
    res = await fetch("/StartTest/API/getTestInfo", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function startATest() {
    res = await fetch("/StartTest/API/startATest", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function stopATest() {
    res = await fetch("/StartTest/API/stopATest", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    return await res.json()
}

async function getTestIp() {
    res = await fetch("/StartTest/API/getTestIp", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })

    curIP = await res.text()
    return curIP
}

async function exportFileCSV() {
    res = await fetch("/StartTest/API/exportCsv", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })

    return await res.json()
}


async function getCandinateList() {
    res = await fetch("/StartTest/API/getCandinateList", { method: 'POST', body: JSON.stringify({ uuid: uuid }) })
    resJ = await res.json()
    return await resJ.candinates
}

function checkCanMark(mark, done) {
    if (done) { return mark } else { return "Chưa nộp bài" }
}

function delay(ms) {
    return new Promise(resolve => {
        setTimeout(resolve, ms);
    });
}

let candinates
async function updateTestCandinate() {
    candinates = await getCandinateList()

    candinateList = document.getElementById("candinateList")

    for (i = 0; i < candinates.length; i++) {
        candinateList.innerHTML += `
            <tr style="cursor: pointer;" onclick="viewTest('${candinates[i].uuid}')">
                <td class="table-cell">${i + 1}</td>
                <td class="table-cell">${candinates[i].name}</td>
                <td class="table-cell">${candinates[i].class}</td>
                <td class="table-cell">${checkCanMark(candinates[i].mark, candinates[i].isDone)}</td>
                <td class="table-cell">${candinates[i].warnTimes}</td>
            </tr>
        `
    }
}

async function viewTest(uuid) {
    if (testInfo.isStarted)
        window.open("http://" + curIP + "/taketest/" + uuid, "blank_")
    else alert("Bắt đầu bài kiểm tra để có thể xem được kết quả")
}


function showElement(id) {
    document.getElementById(id).classList.remove("hidden-element")
}

function hideElement(id) {
    document.getElementById(id).classList.add("hidden-element")
}


let testInfo
async function updateTestInfo() {
    testInfo = await getTestInfo()

    document.getElementById("testName").innerHTML = `<lable>${testInfo.name}</label>`
    document.getElementById("candinate").innerHTML = `<b>Số lượt làm bài:</b> ${testInfo.numberOfCandinate}`
    document.getElementById("testUUID").innerHTML = `<b>Mã đề thi:</b> ${testInfo.uuid}`

    updateTestCandinate()

    if (testInfo.isStarted) {
        ip = await getTestIp()
        document.getElementById("testName").innerHTML = `<lable style="color:lightgreen;">${testInfo.name}</label>`
        document.getElementById("testStatus").innerHTML = "<b>Trạng thái bài kiểm tra:</b> đang được mở"
        document.getElementById("testAddress").innerHTML = `<a style="cursor: pointer;" onclick="window.open('http://${ip}', '_blank');"><b>Làm bài tại:</b> <u style="color: lightblue;">${ip}</u></a> <i id="copyButton" onclick="copyToClipBoard('${ip}')" class="material-icons" style="font-size: 18px;cursor: pointer;margin-top:3px;">content_copy</i>`

        justAButton.innerHTML = "Dừng bài kiểm tra"

        showElement("qrcode")

        var _qrcode = new QRCode("qrcode", {
            text: `http://${ip}`,
            width: 128,
            height: 128,
            colorDark: "#000000",
            colorLight: "#ffffff",
            correctLevel: QRCode.CorrectLevel.H
        });

        justAButton.addEventListener('click', async () => {
            await stopATest()
            location.reload()
        })
    } else {
        document.getElementById("testStatus").innerHTML = "<b>Trạng thái bài kiểm tra:</b> chưa được mở"
        justAButton.addEventListener('click', async () => {
            await startATest()
            location.reload()
        })
    }
}

async function copyToClipBoard(address) {
    navigator.clipboard.writeText(address);

    document.getElementById("copyButton").innerText = "check"
    await delay(1000)
    document.getElementById("copyButton").innerText = "content_copy"
}

let curChart

function deleteCurChart() {
    curChart.destroy()
    hideElement("charPopup")
}

async function createChart(chartElement) {
    let data = new Map()

    for (i = 1; i <= 10; i++) {
        data.set(i, 0)
    }

    for (i = 0; i < candinates.length; i++) {
        if (candinates[i].isDone) {
            if (data.has(Number(candinates[i].mark))) {
                data.set(Number(candinates[i].mark), data.get(Number(candinates[i].mark)) + 1)
            }
            else {
                data.set(Number(candinates[i].mark), 1)
            }
        }
    }

    // sort data
    data = new Map([...data.entries()].sort((a, b) => a[0] - b[0]));

    console.log(data)

    curChart = new Chart(
        chartElement,
        {
            type: 'bar',
            data: {
                labels: Array.from(data.keys()),
                datasets: [
                    {
                        label: 'Số học sinh đạt điểm',
                        data: Array.from(data.values())
                    }
                ]
            }
        }
    );
}