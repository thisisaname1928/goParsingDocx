#include "ui/ui.hpp"
#include "utils.hpp"
#include <QApplication>
#include <QMainWindow>
#include <QVBoxLayout>
#include <QWebEngineView>
#include <QtWebView/QtWebView>
#include <cstring>
#include <filesystem>
#include <iostream>
#include <qcoreapplication.h>
#include <qicon.h>
#include <qlabel.h>
#include <qmainwindow.h>
#include <qnamespace.h>
#include <qwidget.h>

int main(int argc, char *argv[]) {
  QApplication::setHighDpiScaleFactorRoundingPolicy(
      Qt::HighDpiScaleFactorRoundingPolicy::PassThrough);
  QApplication::setAttribute(Qt::AA_UseHighDpiPixmaps);

  QString path = QString::fromStdString(getBinaryPath().string());

  for (int i = path.size() - 1; i >= 0; i--) {
    if (path[i].toLatin1() == '\\' || path[i].toLatin1() == '/') {
      break;
    } else {
      path.remove(i, 1);
    }
  }

  if (argc > 1)
    if (strcmp(argv[1], "dev") == 0) {
      std::cout << "Running under debugmode!\n";
      devMode = true;
    }

  if (devMode) {
    exePath = path.toStdString() + "../..";
    std::filesystem::path iconPath = exePath / "icon.png";
    std::cout << iconPath.string() << '\n';
  } else
    exePath = path.toStdString();

  DouglasApp app(argc, argv);
  std::filesystem::path iconPath = exePath / "icon.png";

  return app.exec();
}