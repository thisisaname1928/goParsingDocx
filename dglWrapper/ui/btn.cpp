#include "ui.hpp"
#include <qboxlayout.h>
#include <qfontdatabase.h>
#include <qlabel.h>
#include <qpushbutton.h>
#include <qwidget.h>

SimpleIconBtn::SimpleIconBtn(QWidget *parent, QString label, QString iconName)
    : QPushButton(parent) {
  QHBoxLayout *layout = new QHBoxLayout(this);
  QLabel *icon = new QLabel(iconName);
  icon->setProperty("class", "SimpleBtnIcon");
  QLabel *l = new QLabel(label);
  l->setProperty("class", "SimpleIconBtnLabel");

  layout->addWidget(icon, 1);
  layout->addWidget(l, 7);

  setProperty("class", "SimpleIconBtn");
  setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Fixed);
}

IconOnlyBtn::IconOnlyBtn(QWidget *parent, QString iconName)
    : QPushButton(parent) {
  setText(iconName);
  setProperty("class", "IconOnlyBtn");
}