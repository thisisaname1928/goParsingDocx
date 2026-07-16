#pragma once

#include "btn.hpp"
#include "home.hpp"
#include <qwidget.h>

class DouglasEditPage : public QWidget {
public:
  DouglasEditPage(QWidget *parent = nullptr);
  IconOnlyBtn *boldBtn;
  IconOnlyBtn *italicBtn;
  IconOnlyBtn *underlineBtn;
  IconOnlyBtn *highlightBtn;
};