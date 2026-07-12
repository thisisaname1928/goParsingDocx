#pragma once

#include <QApplication>
#include <QLabel>
#include <QMainWindow>
#include <QPushButton>
#include <qcoreapplication.h>
#include <qmainwindow.h>
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