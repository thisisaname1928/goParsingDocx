#include "ui.hpp"
#include "../core/load.hpp"
#include "../utils/utils.hpp"
#include "home.hpp"
#include <QGraphicsDropShadowEffect>
#include <QtConcurrent>
#include <cstddef>
#include <filesystem>
#include <iostream>
#include <qapplication.h>
#include <qboxlayout.h>
#include <qcoreapplication.h>
#include <qdebug.h>
#include <qeventloop.h>
#include <qfontdatabase.h>
#include <qframe.h>
#include <qicon.h>
#include <qlabel.h>
#include <qlayout.h>
#include <qmainwindow.h>
#include <qnamespace.h>
#include <qobjectdefs.h>
#include <qpushbutton.h>
#include <qtconcurrentrun.h>
#include <qthread.h>
#include <qtimer.h>
#include <qwidget.h>

QIcon *appIcon;

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

  // load lucide icon
  int fontId = QFontDatabase::addApplicationFont(
      QString::fromStdString((exePath / "lucide.ttf").string()));

  if (fontId == -1) {
    loadingFail("missing lucide icon!");
    return QApplication::exec();
  }

  std::cout
      << QFontDatabase::applicationFontFamilies(fontId).at(0).toStdString()
      << '\n';

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

SideBar::SideBar(QWidget *parent, MainContent *mainContent)
    : QFrame(parent), mainContent(mainContent) {
  setObjectName("sidebar");
  QLabel *label = new QLabel("Side bar here");

  homeTab = new SimpleIconBtn(nullptr, "Trang chủ", "house");
  editTab = new SimpleIconBtn(nullptr, "Soạn", "file-pen-line");
  createTab = new SimpleIconBtn(nullptr, "Tạo bài kiểm tra", "radio");

  homeTab->setCheckable(true);
  editTab->setCheckable(true);
  createTab->setCheckable(true);

  connect(homeTab, &QPushButton::clicked, this, &SideBar::handleHomeTabClick);
  connect(editTab, &QPushButton::clicked, this, &SideBar::handleEditTabClick);
  connect(createTab, &QPushButton::clicked, this,
          &SideBar::handleCreateTabClick);

  handleHomeTabClick();

  QVBoxLayout *layout = new QVBoxLayout(this);
  layout->setAlignment(Qt::AlignTop);
  layout->setContentsMargins(0, 12, 0, 0);
  layout->setSpacing(4);
  layout->addWidget(new sidebarHdr(), 1);
  layout->addWidget(homeTab, 1);
  layout->addWidget(editTab, 1);
  layout->addWidget(createTab, 1);
}

enum SidebarTab { SB_CREATE_TAB = 1, SB_HOME_TAB = 2, SB_EDIT_TAB };

void SideBar::handleCreateTabClick() {
  curTab = SB_CREATE_TAB;

  homeTab->setChecked(false);
  editTab->setChecked(false);
  createTab->setChecked(true);

  mainContent->setMainContentLabel("Tạo bài kiểm tra");
}
void SideBar::handleEditTabClick() {
  curTab = SB_EDIT_TAB;

  homeTab->setChecked(false);
  editTab->setChecked(true);
  createTab->setChecked(false);

  mainContent->setMainContentLabel("Soạn");
}
void SideBar::handleHomeTabClick() {
  curTab = SB_EDIT_TAB;

  homeTab->setChecked(true);
  editTab->setChecked(false);
  createTab->setChecked(false);

  mainContent->setMainContentLabel("Trang chủ");
}

void MainContent::setMainContentLabel(QString title) { label->setText(title); }

void MainContent::setMainContent(DouglasHomePage *page) {
  if (curPage == nullptr) {
    layout->addWidget(page);
    curPage = page;
  } else {
    layout->replaceWidget(curPage, page);
    curPage = page;
  }
}

MainContent::MainContent(QWidget *parent) : QWidget(parent) {
  curPage = nullptr;
  setObjectName("mainContent");
  setAttribute(Qt::WA_StyledBackground, true);
  setContentsMargins(0, 0, 20, 20);
  label = new QLabel("Main content here");
  label->setProperty("class", "MainContentTitle");
  label->setAlignment(Qt::AlignTop);

  layout = new QVBoxLayout(this);
  layout->addWidget(label);
  layout->setSpacing(0);
  layout->setAlignment(Qt::AlignTop);

  QGraphicsDropShadowEffect *shadow = new QGraphicsDropShadowEffect(this);
  shadow->setBlurRadius(24);
  shadow->setXOffset(0);
  shadow->setYOffset(6);
  shadow->setColor(QColor(0, 0, 0, 40));

  setGraphicsEffect(shadow);
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

  MainContent *mainContent = new MainContent(this);
  DouglasHomePage *page = new DouglasHomePage(mainContent);
  mainContent->setMainContent(page);

  layout = new QHBoxLayout(this);
  layout->addWidget(new SideBar(this, mainContent), 1);
  layout->addWidget(mainContent, 5);

  layout->setContentsMargins(0, 10, 10, 10);
  layout->setSpacing(0);
}

DouglasMainWindow::DouglasMainWindow() : QMainWindow() {
  setWindowTitle("Douglas");
  resize(1200, 800);

  mainWidget = new DouglasMainWidget(this);

  setCentralWidget(mainWidget);
}

DouglasMainWindow::~DouglasMainWindow() {}