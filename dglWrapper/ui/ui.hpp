#pragma once

#include "editor.hpp"
#include "home.hpp"
#include <QApplication>
#include <QLabel>
#include <QMainWindow>
#include <QPushButton>
#include <qboxlayout.h>
#include <qcoreapplication.h>
#include <qmainwindow.h>
#include <qobjectdefs.h>
#include <qpushbutton.h>
#include <qwidget.h>

class sidebarHdr : public QWidget {
  Q_OBJECT
public:
  sidebarHdr(QWidget *parent = nullptr);
};

class IconOnlyBtn : public QPushButton {
  Q_OBJECT
public:
  IconOnlyBtn(QWidget *parent, QString icon = "");
};

class SimpleIconBtn : public QPushButton {
  Q_OBJECT
public:
  SimpleIconBtn(QWidget *parent = nullptr, QString label = "",
                QString icon = "");
};

class MainContent : public QWidget {
public:
  MainContent(QWidget *parent = nullptr);
  void setMainContentLabel(QString title);
  void setMainContent(QWidget *page);

private:
  QVBoxLayout *layout;
  QLabel *label;
  QWidget *curPage;
};

class DouglasMainWidget : public QWidget {
public:
  DouglasMainWidget(QWidget *parent = nullptr);
  void switchPage(int page);

private:
  MainContent *mainContent;
  QHBoxLayout *layout;
  DouglasHomePage *homePage;
  DouglasEditPage *editPage;
};

class SideBar : public QFrame {
  Q_OBJECT
public:
  SideBar(DouglasMainWidget *parent = nullptr,
          MainContent *mainContent = nullptr);

public slots:
  void handleHomeTabClick();
  void handleEditTabClick();
  void handleCreateTabClick();

private:
  SimpleIconBtn *homeTab;
  SimpleIconBtn *editTab;
  SimpleIconBtn *createTab;
  int curTab;
  MainContent *mainContent;
};

class DouglasLoadingUI : public QMainWindow {
public:
  DouglasLoadingUI();
  ~DouglasLoadingUI();
  QLabel *label;

private:
  QWidget *mainWidget;
  QLayout *layout;
};

class DouglasMainWindow : public QMainWindow {
public:
  DouglasMainWindow();
  ~DouglasMainWindow();

private:
  QWidget *mainWidget;
};

class DouglasApp : public QApplication {
  Q_OBJECT
public:
  DouglasApp(int argc, char **argv);
  ~DouglasApp();
  int exec();

public slots:
  void doneLoading();
  void loadingFail(QString reason);

private:
  DouglasLoadingUI *loadingWindow;
  DouglasMainWindow *mainWindow;
};

extern QIcon *appIcon;