#include "editor.hpp"
#include "btn.hpp"
#include "ui.hpp"
#include <QAction>
#include <QGraphicsDropShadowEffect>
#include <QList>
#include <QSplitter>
#include <QTextEdit>
#include <qaction.h>
#include <qboxlayout.h>
#include <qbrush.h>
#include <qcolor.h>
#include <qfont.h>
#include <qkeysequence.h>
#include <qlabel.h>
#include <qlist.h>
#include <qnamespace.h>
#include <qpushbutton.h>
#include <qsizepolicy.h>
#include <qsplitter.h>
#include <qtextedit.h>
#include <qtextformat.h>
#include <qwidget.h>

class Toolbar : public QWidget {
public:
  Toolbar(QWidget *parent = nullptr) : QWidget(parent) {
    setObjectName("EditPageToolbar");
    setAttribute(Qt::WA_Hover, true);
  }
};

class TextEditArea : public QTextEdit {
private:
  QWidget *realParent;

public:
  bool isFullyBold() {
    QTextCursor cursor = textCursor();

    if (!cursor.hasSelection())
      return cursor.charFormat().fontWeight() == QFont::Bold;

    int start = cursor.selectionStart();
    int end = cursor.selectionEnd();

    QTextCursor c(cursor.document());
    c.setPosition(start);
    while (c.position() < end) {
      c.movePosition(QTextCursor::NextCharacter, QTextCursor::KeepAnchor);
      if (c.charFormat().fontWeight() != QFont::Bold)
        return false;
      c.setPosition(c.position());
    }
    return true;
  }

  bool isFullyItalic() {
    QTextCursor cursor = textCursor();

    if (!cursor.hasSelection())
      return cursor.charFormat().fontItalic();

    int start = cursor.selectionStart();
    int end = cursor.selectionEnd();

    QTextCursor c(cursor.document());
    c.setPosition(start);
    while (c.position() < end) {
      c.movePosition(QTextCursor::NextCharacter, QTextCursor::KeepAnchor);
      if (!c.charFormat().fontItalic())
        return false;
      c.setPosition(c.position());
    }
    return true;
  }

  bool isFullyUnderline() {
    QTextCursor cursor = textCursor();

    if (!cursor.hasSelection())
      return cursor.charFormat().fontUnderline();

    int start = cursor.selectionStart();
    int end = cursor.selectionEnd();

    QTextCursor c(cursor.document());
    c.setPosition(start);
    while (c.position() < end) {
      c.movePosition(QTextCursor::NextCharacter, QTextCursor::KeepAnchor);
      if (!c.charFormat().fontUnderline())
        return false;
      c.setPosition(c.position());
    }
    return true;
  }

  bool isFullyHighlight() {
    QTextCursor cursor = textCursor();

    if (!cursor.hasSelection())
      return cursor.charFormat().background().style() != Qt::NoBrush;

    int start = cursor.selectionStart();
    int end = cursor.selectionEnd();

    QTextCursor c(cursor.document());
    c.setPosition(start);
    while (c.position() < end) {
      c.movePosition(QTextCursor::NextCharacter, QTextCursor::KeepAnchor);
      if (c.charFormat().background().style() == Qt::NoBrush)
        return false;
      c.setPosition(c.position());
    }
    return true;
  }

  void onSelecting() {
    if (isFullyBold()) {
      editPage()->boldBtn->setChecked(true);
    } else {
      editPage()->boldBtn->setChecked(false);
    }

    if (isFullyItalic()) {
      editPage()->italicBtn->setChecked(true);
    } else {
      editPage()->italicBtn->setChecked(false);
    }

    if (isFullyUnderline()) {
      editPage()->underlineBtn->setChecked(true);
    } else {
      editPage()->underlineBtn->setChecked(false);
    }

    if (isFullyHighlight()) {
      editPage()->highlightBtn->setChecked(true);
    } else {
      editPage()->highlightBtn->setChecked(false);
    }
  }

  DouglasEditPage *editPage() {
    return static_cast<DouglasEditPage *>(realParent);
  }

  TextEditArea(QWidget *parent = nullptr)
      : QTextEdit(parent), realParent(parent) {
    setObjectName("TextEditArea");

    document()->setDocumentMargin(14);
    setPlaceholderText("Câu 1: Nhập câu hỏi...\nA. Đáp án A\nB. Đáp án B\nC. "
                       "Đáp án C\nD. Đáp án D");

    connect(this, &TextEditArea::selectionChanged, this,
            &TextEditArea::onSelecting);
    connect(this, &QTextEdit::cursorPositionChanged, this,
            &TextEditArea::onSelecting);
  }

  bool isSelecting() { return textCursor().hasSelection(); }

  void toggleBold() {
    if (textCursor().hasSelection()) {
      if (isFullyBold()) {
        QTextCharFormat tmp;
        tmp.setFontWeight(QFont::Normal);
        textCursor().mergeCharFormat(tmp);
        editPage()->boldBtn->setChecked(true);
      } else {
        QTextCharFormat tmp;
        tmp.setFontWeight(QFont::Bold);
        textCursor().mergeCharFormat(tmp);

        editPage()->boldBtn->setChecked(false);
      }
    } else {
      if (isFullyBold()) {
        QTextCharFormat tmp;
        tmp.setFontWeight(QFont::Normal);

        mergeCurrentCharFormat(tmp);
      } else {
        QTextCharFormat tmp;
        tmp.setFontWeight(QFont::Bold);

        mergeCurrentCharFormat(tmp);
      }
    }

    // QTextCursor cursor = this->textCursor();

    // QTextCharFormat format;

    // format.setFontWeight(QFont::Bold);
    // cursor.mergeCharFormat(format);
    // mergeCurrentCharFormat(format);
  }

  void toggleItalic() {
    if (isFullyItalic()) {
      if (textCursor().hasSelection()) {
        QTextCharFormat format;
        format.setFontItalic(false);
        textCursor().mergeCharFormat(format);
      } else {
        QTextCharFormat format;
        format.setFontItalic(false);
        mergeCurrentCharFormat(format);
      }
    } else {
      if (textCursor().hasSelection()) {
        QTextCharFormat format;
        format.setFontItalic(true);
        textCursor().mergeCharFormat(format);
      } else {
        QTextCharFormat format;
        format.setFontItalic(true);
        mergeCurrentCharFormat(format);
      }
    }
  }

  void toggleUnderline() {
    if (isFullyUnderline()) {
      if (textCursor().hasSelection()) {
        QTextCharFormat format;
        format.setFontUnderline(false);
        textCursor().mergeCharFormat(format);
      } else {
        QTextCharFormat format;
        format.setFontUnderline(false);
        mergeCurrentCharFormat(format);
      }
    } else {
      if (textCursor().hasSelection()) {
        QTextCharFormat format;
        format.setFontUnderline(true);
        textCursor().mergeCharFormat(format);
      } else {
        QTextCharFormat format;
        format.setFontUnderline(true);
        mergeCurrentCharFormat(format);
      }
    }
  }

  void toggleHighlight() {
    if (isFullyHighlight()) {
      if (textCursor().hasSelection()) {
        QTextCharFormat format = textCursor().charFormat();
        format.clearProperty(QTextFormat::BackgroundBrush);
        textCursor().setCharFormat(format);
      } else {
        QTextCharFormat format = textCursor().charFormat();
        format.clearProperty(QTextFormat::BackgroundBrush);
        setCurrentCharFormat(format);
      }
    } else {
      if (textCursor().hasSelection()) {
        QTextCharFormat format;
        format.setBackground(QColor("yellow"));
        textCursor().mergeCharFormat(format);
      } else {
        QTextCharFormat format;
        format.setBackground(QColor("yellow"));
        mergeCurrentCharFormat(format);
      }
    }
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
  QHBoxLayout *toolbarLayout = new QHBoxLayout(toolbar);

  toolbarLayout->setContentsMargins(10, 7, 10, 7);
  toolbarLayout->setSpacing(4);

  IconOnlyBtn *fileBtn = new IconOnlyBtn(toolbar, "folder");

  boldBtn = new IconOnlyBtn(toolbar, "bold");
  boldBtn->setCheckable(true);
  boldBtn->setFocusPolicy(Qt::NoFocus);

  italicBtn = new IconOnlyBtn(toolbar, "italic");
  italicBtn->setCheckable(true);
  italicBtn->setFocusPolicy(Qt::NoFocus);

  underlineBtn = new IconOnlyBtn(toolbar, "underline");
  underlineBtn->setCheckable(true);
  underlineBtn->setFocusPolicy(Qt::NoFocus);

  highlightBtn = new IconOnlyBtn(toolbar, "highlighter");
  highlightBtn->setCheckable(true);
  highlightBtn->setFocusPolicy(Qt::NoFocus);

  IconOnlyBtn *undoBtn = new IconOnlyBtn(toolbar, "undo-2");
  undoBtn->setShortcut(QKeySequence::Undo);
  undoBtn->setFocusPolicy(Qt::NoFocus);

  IconOnlyBtn *redoBtn = new IconOnlyBtn(toolbar, "redo-2");
  redoBtn->setShortcut(QKeySequence::Redo);
  redoBtn->setFocusPolicy(Qt::NoFocus);

  SimpleIconBtn *addQuestionBtn = new SimpleIconBtn(toolbar, "Câu hỏi", "plus");

  toolbarLayout->addWidget(fileBtn);
  toolbarLayout->addWidget(boldBtn);
  toolbarLayout->addWidget(italicBtn);
  toolbarLayout->addWidget(underlineBtn);
  toolbarLayout->addWidget(highlightBtn);
  toolbarLayout->addWidget(undoBtn);
  toolbarLayout->addWidget(redoBtn);
  toolbarLayout->addWidget(addQuestionBtn);
  toolbarLayout->addStretch(1);

  // tool bar trigger
  connect(boldBtn, &QPushButton::pressed, editArea, &TextEditArea::toggleBold);
  connect(italicBtn, &QPushButton::pressed, editArea,
          &TextEditArea::toggleItalic);
  connect(underlineBtn, &QPushButton::pressed, editArea,
          &TextEditArea::toggleUnderline);
  connect(highlightBtn, &QPushButton::pressed, editArea,
          &TextEditArea::toggleHighlight);
  connect(undoBtn, &QPushButton::pressed, editArea, &TextEditArea::undo);
  connect(redoBtn, &QPushButton::pressed, editArea, &TextEditArea::redo);

  toolbar->setLayout(toolbarLayout);

  layout->addWidget(splitter, 1);
}