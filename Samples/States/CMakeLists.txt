cmake_minimum_required(VERSION 3.7)
project(States)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -Wall")

include_directories( ${WINDOW_SOURCE_DIR}../../dep/include )
link_directories( ${WINDOW_SOURCE_DIR}../../dep/lib/Linux )

include_directories( ${WINDOW_SOURCE_DIR}../../Engine/include )
link_directories( ${WINDOW_SOURCE_DIR}../../cmake-build-debug/Engine )

# DEBUG|RELEASE|RELWITHDEBINFO|MINSIZEREL
set(CMAKE_BUILD_TYPE DEBUG)

add_executable(States main.cpp)

target_link_libraries( States Paranoia )