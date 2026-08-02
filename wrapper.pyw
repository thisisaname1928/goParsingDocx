import webview
import zipfile
import time
import subprocess
import sys
import requests
import json

host = "localhost"
port = 8080

try:
    webview.platforms.gtk.BrowserView.settings.set_property('enable-developer-extras', True)
    webview.platforms.gtk.BrowserView.settings.set_property('enable-local-storage', True)
    webview.platforms.gtk.BrowserView.settings.set_property('enable-dom-storage', True)
except:
    pass

def update(path):
    try:
        f = zipfile.ZipFile(path, 'r')
        flist = f.namelist()
        for i in flist:
            f.extract(i, path="./")
        return True
    except:
        return False
    
def check4Update():
    try:
        f = open("./appVersion.json", 'r', encoding='utf-8')
        dat = json.load(f)

        if dat["shouldUpdate"]:
            print("updatingg...")
            r = update('update.zip')
            if r:
                dat["shouldUpdate"] = False
                r = json.dumps(dat)
                f = open("./appVersion.json", "w", encoding="utf-8")
                f.write(r)
                f.close()
                print("ok")
            else:
                print("not ok")
    except:
        print("not ok")

class Api:
    def __init__(self):
        self._window = None

    def set_window(self, window):
        self._window = window

    def select_file(self):
        if not self._window:
            return ""
        file_types = ('Docx Files (*.docx)', 'All files (*.*)')
        result = self._window.create_file_dialog(webview.OPEN_DIALOG, allow_multiple=False, file_types=file_types)
        if result and len(result) > 0:
            return result[0]
        return ""

    def save_file(self, filename="exported.dou"):
        if not self._window:
            return ""
        file_types = ('Douglas Files (*.dou)', 'All files (*.*)')
        result = self._window.create_file_dialog(webview.SAVE_DIALOG, save_filename=filename, file_types=file_types)
        if result:
            if isinstance(result, (list, tuple)) and len(result) > 0:
                return result[0]
            return str(result)
        return ""

import os

if __name__ == '__main__':
    check4Update()
    
    executable = "goParsingDocx.exe" if sys.platform == "win32" else "./goParsingDocx"
    proccess = subprocess.Popen([executable])

    while True:
        try:
            response = requests.get(f'http://{host}:{port}/check')
            if response.status_code == 200:
                break
        except:
            time.sleep(0.5)
    
    webview.settings['ALLOW_DOWNLOADS'] = True
    api = Api()
    window = webview.create_window(title="Douglas", url=f'http://{host}:{port}/Home', width=1424, height=700, js_api=api)
    api.set_window(window)

    webview.start(debug=True, private_mode=True, storage_path=os.path.abspath("cache"))

    if sys.platform == "win32":
        proccess.terminate()
    else:
        proccess.kill()

    proccess.wait()