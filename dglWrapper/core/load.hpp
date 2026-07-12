#pragma once

#include <QCoreApplication>
#include <qobject.h>
#include <qobjectdefs.h>

class Loader : public QObject {
  Q_OBJECT
public:
  Loader();
  void start();
signals:
  void done();
  void fail(QString reason);
};
