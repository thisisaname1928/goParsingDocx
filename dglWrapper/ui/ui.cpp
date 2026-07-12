#include "ui.hpp"
#include "../core/load.hpp"
#include "../utils/utils.hpp"
#include <QtConcurrent>
#include <cstddef>
#include <filesystem>
#include <iostream>
#include <qapplication.h>
#include <qboxlayout.h>
#include <qcoreapplication.h>
#include <qdebug.h>
#include <qeventloop.h>
#include <qframe.h>
#include <qicon.h>
#include <qlabel.h>
#include <qlayout.h>
#include <qmainwindow.h>
#include <qnamespace.h>
#include <qobjectdefs.h>
#include <qtconcurrentrun.h>
#include <qthread.h>
#include <qtimer.h>
#include <qwidget.h>

DouglasLoadingUI::DouglasLoadingUI() : QMainWindow() {
  setWindowTitle("Douglas");
  resize(600, 400);
  setWindowFlags(Qt::FramelessWindowHint);
  setAttribute(Qt::WA_DeleteOnClose);
  // setAttribute(Qt::WA_TranslucentBackground);

  mainWidget = new QWidget(this);
  layout = new QVBoxLayout(mainWidget);
  label = new QLabel("Douglas is loading...");
  label->setAlignment(Qt::AlignBottom);
  layout->addWidget(label);

  setCentralWidget(mainWidget);
}

DouglasLoadingUI::~DouglasLoadingUI() {}

int DouglasApp::exec() {
  // setup the loading window
  loadingWindow = new DouglasLoadingUI();
  loadingWindow->setWindowIcon(*appIcon);

  loadingWindow->show();

  Loader *loader = new Loader;
  connect(loader, &Loader::done, this, &DouglasApp::doneLoading);
  connect(loader, &Loader::fail, this, &DouglasApp::loadingFail);

  // loadup styles

  std::string qssSrc;
  std::cout << (exePath / "styles.qss").string() << '\n';
  try {
    qssSrc = readFileContent(exePath / "styles.qss");
  } catch (...) {
    loadingFail("can't open styles.qss!");
    return QApplication::exec();
  }

  setStyleSheet(QString::fromStdString(qssSrc));

  // init main window
  mainWindow = new DouglasMainWindow();
  mainWindow->setWindowIcon(*appIcon);

  QtConcurrent::run([this, loader]() {
    // let loader control
    loader->start();
  });

  return QApplication::exec();
}

void DouglasApp::doneLoading() {
  loadingWindow->close();
  mainWindow->show();
}
void DouglasApp::loadingFail(QString reason) {
  loadingWindow->label->setText(
      QString::fromStdString("Loading failed! Reason: ") + reason);
}

DouglasApp::DouglasApp(int argc, char **argv) : QApplication(argc, argv) {
  // load the icon
  appIcon = new QIcon(QString::fromStdString(
      std::filesystem::path(exePath / "icon.png").string()));

  if (!appIcon->isNull()) {
    setWindowIcon(*appIcon);
  } else {
    std::cout << "Douglas Icon not found!\n";
  }

  std::cout << "Douglas Wrapper ok!\n";
}

DouglasApp::~DouglasApp() {
  if (appIcon != nullptr) {
    delete appIcon;
  }
}

class SideBar : public QFrame {
public:
  SideBar(QWidget *parent = nullptr);
};

SideBar::SideBar(QWidget *parent) : QFrame(parent) {
  setObjectName("sidebar");
  QLabel *label = new QLabel("Side bar here");

  QVBoxLayout *layout = new QVBoxLayout(this);
  layout->setContentsMargins(0, 0, 0, 0);
  layout->setSpacing(0);
  layout->addWidget(label);
}

class MainContent : public QWidget {
public:
  MainContent();
};

MainContent::MainContent() : QWidget() {
  setObjectName("mainContent");
  setAttribute(Qt::WA_StyledBackground, true);
  QLabel *label = new QLabel("Main content here");

  QVBoxLayout *layout = new QVBoxLayout(this);
  layout->setContentsMargins(0, 0, 0, 0);
  layout->addWidget(label);
  layout->setSpacing(0);
}

class DouglasMainWidget : public QWidget {
public:
  DouglasMainWidget(QWidget *parent = nullptr);

private:
  QHBoxLayout *layout;
};

DouglasMainWidget::DouglasMainWidget(QWidget *parent) : QWidget(parent) {
  setAttribute(Qt::WA_StyledBackground, true);
  setObjectName("mainWidget");

  layout = new QHBoxLayout(this);
  layout->addWidget(new SideBar(this), 1);
  layout->addWidget(new MainContent, 5);

  layout->setContentsMargins(0, 0, 0, 0);
  layout->setSpacing(0);
}

DouglasMainWindow::DouglasMainWindow() : QMainWindow() {
  setWindowTitle("Douglas");
  resize(600, 400);

  mainWidget = new DouglasMainWidget(this);

  setCentralWidget(mainWidget);
}

DouglasMainWindow::~DouglasMainWindow() {}