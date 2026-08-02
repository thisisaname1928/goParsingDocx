const modeToggleButton = document.getElementById('modeToggleButton');
const msg = document.getElementById('msg')
const configBox = document.getElementById('configBox')
const exportButton = document.getElementById('exp')
const body = document.body;
function getUUIDFromURL() {
    const url = new URL(window.location.href);
    const pathParts = url.pathname.split('/').filter(p => p.length > 0);
    const uuidIndex = pathParts.indexOf('UUID');
    if (uuidIndex !== -1 && uuidIndex + 1 < pathParts.length) {
        return pathParts[uuidIndex + 1].split('?')[0];
    }
    return "";
}
let UUID = getUUIDFromURL();
const parsedUrl = new URL(location.href);
const parseExportType = parsedUrl.searchParams.get("exportType") || "useDocx";


document.addEventListener('DOMContentLoaded', () => {

    UUID = getUUIDFromURL();

    exportButton.addEventListener('click', async () => {
        expData = getConfig()
        if (!expData.status) {
            alert(expData.msg)
            return
        }

        await exportTest(expData)
    })

    initConf()
})

function renderStypeConfig(quesType, maxQues) {
    return `<div class="config-container">
                <label class="config-label">Câu loại ${quesType}, số câu ${maxQues}:</label>
                
                <table>
                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Số câu:</label>
                    </td>
                    <td>
                        <input id="${quesType}.N" type="number" class="config-input" placeholder="Số câu mỗi đề">
                    </td>
                </tr>
                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Số điểm:</label>
                    </td>
                    <td>
                        <input id="${quesType}.Point" type="number" class="config-input" placeholder="Số điểm mỗi câu">
                    </td>
                </tr>
                </table>
            </div>`
}

let confObj = {};

async function initConf() {
    res = await fetch4ExportConfig()
    obj = await res.json()

    confObj = obj

    if (!obj.status) {
        if (obj.msg == "invalid docx file") {
            msg.innerText = "File được xuất không hợp lệ!"
        } else {
            msg.innerText = obj.msg
        }
        return
    }

    configBox.innerHTML += `<div class="config-container">
                <div class="input-config-row">
                    <label class="config-label">Tác giả:</label>
                </div>
                <div class="input-config-row">
                    <input id="author" type="text" class="config-input" placeholder="Tên người ra đề">
                </div>
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Mật khẩu cho đề:</label>
                <input id="key" type="text" class="config-input" placeholder="Để trống nếu không dùng mã hóa">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Thời gian làm bài (phút):</label>
                <input id="testDuration" type="number" class="config-input" placeholder="Đặt 0 nếu không có thời gian cố định">
            </div>`

    configBox.innerHTML += `<div class="config-container">
                <label class="config-label">Điểm đúng sai: </label>
                
                <table>
                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Đúng một câu (%):</label>
                    </td>
                    <td>
                        <input id="TNDSPer1" type="number" class="config-input" placeholder="% điểm của câu" value="25">
                    </td>
                </tr>
                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Đúng hai câu (%):</label>
                    </td>
                    <td>
                        <input id="TNDSPer2" type="number" class="config-input" placeholder="% điểm của câu" value="50">
                    </td>
                </tr>
                                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Đúng ba câu (%):</label>
                    </td>
                    <td>
                        <input id="TNDSPer3" type="number" class="config-input" placeholder="% điểm của câu" value="75">
                    </td>
                </tr>
                                <tr class="input-config-row">
                    <td>
                        <label  class="sub-config-label">Đúng bốn câu (%):</label>
                    </td>
                    <td>
                        <input id="TNDSPer4" type="number" class="config-input" placeholder="% điểm của câu" value="100">
                    </td>
                </tr>
                </table>
            </div>`

    for (i = 0; i < obj.stype.length; i++) {
        configBox.innerHTML += renderStypeConfig(obj.stype[i].stype, obj.stype[i].N)
    }

    alert("ok")
}

function getTNDSPointCalcInfo() {
    TNDSPer1 = document.getElementById('TNDSPer1').value
    TNDSPer2 = document.getElementById('TNDSPer2').value
    TNDSPer3 = document.getElementById('TNDSPer3').value
    TNDSPer4 = document.getElementById('TNDSPer4').value

    return [Number(TNDSPer1), Number(TNDSPer2), Number(TNDSPer3), Number(TNDSPer4)]
}

function getConfig() {
    author = document.getElementById("author").value
    key = document.getElementById("key").value
    testDuration = Math.abs(Number(document.getElementById("testDuration").value)) // abs to prevent user do somethings wrong=)))
    if (Number.isNaN(testDuration)) testDuration = 0;
    stype = []

    if (key.length > 16) { return { status: false, UUID: UUID, testDuration: testDuration, msg: "Mật khẩu quá dài", author: author, key: key, stype: stype } }


    for (i = 0; i < confObj.stype.length; i++) {
        numberOfQuesPerTest = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.N`).value))
        pointPerQues = Math.abs(Number(document.getElementById(`${confObj.stype[i].stype}.Point`).value))

        if (Number.isNaN(numberOfQuesPerTest) || Number.isNaN(pointPerQues)) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `Dữ liệu nhập cho loại câu ${confObj.stype[i].stype} không hợp lệ!`, author: author, key: key, stype: stype }
        }

        if (numberOfQuesPerTest > confObj.stype[i].N) {
            return { status: false, UUID: UUID, testDuration: testDuration, msg: `Số câu mỗi đề của loại ${confObj.stype[i].stype} lớn hơn số câu tồn tại`, author: author, key: key, stype: stype }
        }

        stype.push({ stype: confObj.stype[i].stype, N: numberOfQuesPerTest, Point: pointPerQues })
    }

    return { status: true, UUID: UUID, testDuration: testDuration, msg: "ok", author: author, key: key, stype: stype, exportType: parseExportType, additionalExportData: { TNDSPointCalcInfo: getTNDSPointCalcInfo() } }
}

async function exportTest(obj) {
    res = await fetch("/Export/API/export", { method: "POST", body: JSON.stringify(obj) })
    dat = await res.json()
    if (dat.status == false) {
        alert(dat.msg)
        return
    }

    if (window.pywebview && window.pywebview.api && window.pywebview.api.save_file) {
        try {
            const savePath = await window.pywebview.api.save_file("exported.dou");
            if (savePath) {
                const saveRes = await fetch("/Export/API/saveToPath", {
                    method: "POST",
                    body: JSON.stringify({ UUID: UUID, targetPath: savePath })
                });
                const saveJson = await saveRes.json();
                if (saveJson.status) {
                    alert("Xuất đề thành công! Đã lưu tại: " + savePath);
                } else {
                    alert("Lỗi khi lưu file: " + saveJson.msg);
                }
            }
        } catch (e) {
            console.error("Error with OS save file dialog:", e);
            alert('Xuất đề thành công!');
            window.location.href = "/Export/Download/UUID/" + UUID + ".dou";
        }
    } else {
        alert('Xuất đề thành công!')
        const link = document.createElement("a")
        link.href = "/Export/Download/UUID/" + UUID + ".dou"
        link.download = 'exported.dou'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    }
}

function fetch4ExportConfig() {
    res = fetch("/Export/API/getConfig", { method: "POST", body: JSON.stringify({ UUID: UUID, exportType: parseExportType }) })
    return res
}

