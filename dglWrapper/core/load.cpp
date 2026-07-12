#include "load.hpp"
#include <chrono>
#include <qobjectdefs.h>
#include <thread>

Loader::Loader() : QObject() {}

void Loader::start() {
  // std::this_thread::sleep_for(std::chrono::seconds(2));

  emit done();
}