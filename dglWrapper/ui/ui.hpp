#pragma once

#include <QApplication>
#include <QLabel>
#include <QMainWindow>
#include <qcoreapplication.h>
#include <qmainwindow.h>
#include <qwidget.h>

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
  QIcon *appIcon;
  DouglasLoadingUI *loadingWindow;
  DouglasMainWindow *mainWindow;
};