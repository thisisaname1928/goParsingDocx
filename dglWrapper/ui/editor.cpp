#include "editor.hpp"
#include "ui.hpp"
#include <QGraphicsDropShadowEffect>
#include <QSplitter>
#include <QTextEdit>
#include <qboxlayout.h>
#include <qlabel.h>
#include <qsizepolicy.h>
#include <qsplitter.h>
#include <qtextedit.h>
#include <qwidget.h>

class Toolbar : public QWidget {
public:
  Toolbar(QWidget *parent = nullptr) : QWidget(parent) {
    setObjectName("EditPageToolbar");
    setAttribute(Qt::WA_Hover, true);

    QHBoxLayout *layout = new QHBoxLayout(this);

    layout->setContentsMargins(10, 7, 10, 7);
    layout->setSpacing(4);

    IconOnlyBtn *fileBtn = new IconOnlyBtn(this, "folder");
    IconOnlyBtn *boldBtn = new IconOnlyBtn(this, "bold");
    IconOnlyBtn *italicBtn = new IconOnlyBtn(this, "italic");
    IconOnlyBtn *underlineBtn = new IconOnlyBtn(this, "underline");
    IconOnlyBtn *listBtn = new IconOnlyBtn(this, "list");
    SimpleIconBtn *addQuestionBtn = new SimpleIconBtn(this, "Câu hỏi", "plus");

    layout->addWidget(fileBtn);
    layout->addWidget(boldBtn);
    layout->addWidget(italicBtn);
    layout->addWidget(underlineBtn);
    layout->addWidget(listBtn);
    layout->addWidget(addQuestionBtn);
    layout->addStretch(1);

    setLayout(layout);
  }
};

class TextEditArea : public QTextEdit {
public:
  TextEditArea(QWidget *parent = nullptr) : QTextEdit(parent) {
    setObjectName("TextEditArea");

    document()->setDocumentMargin(14);
    setPlaceholderText("Câu 1: Nhập câu hỏi...\nA. Đáp án A\nB. Đáp án B\nC. "
                       "Đáp án C\nD. Đáp án D");
  }
};

class PreviewArea : public QWidget {
public:
  PreviewArea(QWidget *parent = nullptr) : QWidget(parent) {
    setObjectName("PreviewArea");
  }
};

DouglasEditPage::DouglasEditPage(QWidget *parent) : QWidget(parent) {
  QVBoxLayout *layout = new QVBoxLayout(this);
  setLayout(layout);

  TextEditArea *editArea = new TextEditArea(this);
  PreviewArea *previewArea = new PreviewArea(this);

  editArea->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
  previewArea->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);

  QSplitter *splitter = new QSplitter(this);
  splitter->setChildrenCollapsible(false);
  splitter->setHandleWidth(1);
  splitter->addWidget(editArea);
  splitter->addWidget(previewArea);
  splitter->setStretchFactor(0, 1);
  splitter->setStretchFactor(1, 1);
  splitter->setSizes(QList<int>() << 10000 << 10000);

  Toolbar *toolbar = new Toolbar(this);

  layout->addWidget(toolbar);
  layout->addWidget(splitter, 1);
}