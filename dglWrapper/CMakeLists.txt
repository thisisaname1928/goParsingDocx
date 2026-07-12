cmake_minimum_required(VERSION 3.16)

project(douglas LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

find_package(Qt5 REQUIRED COMPONENTS Core Gui Widgets WebEngineWidgets WebView Concurrent)

if(CMAKE_SOURCE_DIR STREQUAL CMAKE_BINARY_DIR)
    message(FATAL_ERROR "pls do not building this in-source")
endif()

file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS
    ${CMAKE_SOURCE_DIR}/*.cpp
)

file(GLOB_RECURSE HEADERS CONFIGURE_DEPENDS
    ${CMAKE_SOURCE_DIR}/*.h
    ${CMAKE_SOURCE_DIR}/*.hpp
)

file(GLOB_RECURSE UI_FILES CONFIGURE_DEPENDS
    ${CMAKE_SOURCE_DIR}/*.ui
)

list(FILTER SOURCES EXCLUDE REGEX "${CMAKE_BINARY_DIR}")
list(FILTER HEADERS EXCLUDE REGEX "${CMAKE_BINARY_DIR}")
list(FILTER UI_FILES EXCLUDE REGEX "${CMAKE_BINARY_DIR}")

add_executable(${PROJECT_NAME}
    ${SOURCES}
    ${HEADERS}
    ${UI_FILES}
)

foreach(_header ${HEADERS})
    get_filename_component(_dir ${_header} DIRECTORY)
    list(APPEND INCLUDE_DIRS ${_dir})
endforeach()
list(REMOVE_DUPLICATES INCLUDE_DIRS)
target_include_directories(${PROJECT_NAME} PRIVATE ${INCLUDE_DIRS})

target_link_libraries(${PROJECT_NAME}
    PRIVATE
    Qt5::Core
    Qt5::Gui
    Qt5::Widgets
    Qt5::WebEngineWidgets
    Qt5::WebView
    Qt5::Concurrent
)