cmake_minimum_required(VERSION 3.7)
project(Paranoia)

set( MAJOR_VERSION 0.1.0 )
set( MINOR_VERSION 0.1 )

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -Wall")

set(SRC_DIR ${PARANOIA_SOURCE_DIR}/src)
include_directories( ${PARANOIA_SOURCE_DIR}/include )

include_directories( ${PARANOIA_SOURCE_DIR}../dep/include )
link_directories( ${PARANOIA_SOURCE_DIR}../dep/lib/Linux )


# find OpenGL
find_package(OpenGL REQUIRED)
include_directories(${OPENGL_INCLUDE_DIR})
set(ADDITIONAL_LIBRARIES ${OPENGL_LIBRARIES})

# DEBUG|RELEASE|RELWITHDEBINFO|MINSIZEREL
set(CMAKE_BUILD_TYPE DEBUG)


file(GLOB_RECURSE SOURCES *.cpp)
file(GLOB_RECURSE HEADERS "*.h")


set(INSTALL_PATH ${PARANOIA_SOURCE_DIR}/../bin)


add_library(Paranoia SHARED ${SOURCES} ${HEADERS})


target_link_libraries(Paranoia sfml-window sfml-system sfml-graphics z ${ADDITIONAL_LIBRARIES})

set_target_properties( Paranoia PROPERTIES 
                              VERSION ${MAJOR_VERSION} 
                              SOVERSION ${MINOR_VERSION} )
                              
install (TARGETS Paranoia DESTINATION ${INSTALL_PATH})
