cmake_minimum_required(VERSION 3.7)
#project(Paranoia)

set(CMAKE_CXX_STANDARD 11)

add_subdirectory(Engine)

#Samples build
add_subdirectory(Samples/Window)
add_subdirectory(Samples/Threads)