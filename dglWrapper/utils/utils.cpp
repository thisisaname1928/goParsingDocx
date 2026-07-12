#include "utils.hpp"
#include <filesystem>
#include <fstream>
#include <stdexcept>
#if defined(_WIN32)
#include <windows.h>
#elif defined(__linux__)
#include <limits.h>
#include <unistd.h>
#elif defined(__APPLE__)
#include <mach-o/dyld.h>
#include <vector>
#endif

#include <string>

std::filesystem::path getBinaryPath() {
#if defined(_WIN32)
  char buffer[MAX_PATH];
  GetModuleFileNameA(NULL, buffer, MAX_PATH);
  return std::filesystem::path(buffer);
#elif defined(__linux__)
  char buffer[PATH_MAX];
  ssize_t len = readlink("/proc/self/exe", buffer, sizeof(buffer) - 1);
  if (len != -1) {
    buffer[len] = '\0';
    return std::filesystem::path(buffer);
  }
#elif defined(__APPLE__)
  uint32_t size = 0;
  _NSGetExecutablePath(nullptr, &size);
  std::vector<char> buffer(size);
  if (_NSGetExecutablePath(buffer.data(), &size) == 0) {
    return std::filesystem::path(buffer.data());
  }
#endif
  return "";
}

bool devMode = false;
std::filesystem::path exePath;

std::string readFileContent(const std::filesystem::path &path) {
  std::ifstream file(path);
  if (!file.is_open()) {
    throw std::runtime_error("CAN'T OPEN FILE");
  }
  std::stringstream buffer;
  buffer << file.rdbuf();

  return buffer.str();
}