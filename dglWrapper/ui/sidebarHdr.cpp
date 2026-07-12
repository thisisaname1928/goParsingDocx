#include "ui.hpp"
#include <qboxlayout.h>
#include <qlabel.h>
#include <qsize.h>
#include <qwidget.h>

sidebarHdr::sidebarHdr(QWidget *parent) : QWidget(parent) {
  QHBoxLayout *layout = new QHBoxLayout(this);

  QLabel *ico = new QLabel();
  QLabel *label = new QLabel("Douglas");

  label->setProperty("class", "HdrTxt");

  ico->setPixmap(appIcon->pixmap(QSize(22, 22)));
  ico->setFixedSize(22, 22);
  ico->setProperty("class", "HdrIcon");

  layout->addWidget(ico);

  layout->addWidget(label);
  layout->setSpacing(8);
  layout->setContentsMargins(20, 0, 0, 40);
  setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Fixed);

  setProperty("class", "SidebarHdr");
}