#pragma once

#include "home.hpp"
#include <QApplication>
#include <QLabel>
#include <QMainWindow>
#include <QPushButton>
#include <qboxlayout.h>
#include <qcoreapplication.h>
#include <qmainwindow.h>
#include <qobjectdefs.h>
#include <qwidget.h>

class sidebarHdr : public QWidget {
public:
  sidebarHdr(QWidget *parent = nullptr);
};

class SimpleIconBtn : public QPushButton {
public:
  SimpleIconBtn(QWidget *parent = nullptr, QString label = "",
                QString icon = "");
};

class MainContent : public QWidget {
public:
  MainContent(QWidget *parent = nullptr);
  void setMainContentLabel(QString title);
  void setMainContent(DouglasHomePage *page);

private:
  QVBoxLayout *layout;
  QLabel *label;
  DouglasHomePage *curPage;
};

class SideBar : public QFrame {
  Q_OBJECT
public:
  SideBar(QWidget *parent = nullptr, MainContent *mainContent = nullptr);

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