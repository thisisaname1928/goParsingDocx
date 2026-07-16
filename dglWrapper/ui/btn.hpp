#pragma once
#include <QPushButton>

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