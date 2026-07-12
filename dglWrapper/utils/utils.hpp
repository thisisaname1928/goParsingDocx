#pragma once
#include <filesystem>
#include <string>

std::filesystem::path getBinaryPath();

std::string readFileContent(const std::filesystem::path &path);

extern bool devMode;
extern std::filesystem::path exePath;