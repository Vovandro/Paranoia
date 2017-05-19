cmake_minimum_required(VERSION 3.7)
#project(Paranoia)

set(CMAKE_CXX_STANDARD 11)

add_subdirectory(Engine)


set(BUILD_EXAMPLES ON)

if (${BUILD_EXAMPLES})
    #Samples build
    add_subdirectory(Samples/Window)
    add_subdirectory(Samples/Threads)
    add_subdirectory(Samples/Files)
endif ()