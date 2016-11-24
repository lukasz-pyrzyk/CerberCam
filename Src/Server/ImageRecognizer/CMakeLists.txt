cmake_minimum_required(VERSION 3.2 FATAL_ERROR)
project (TensorflowPlayground VERSION 0.1 LANGUAGES CXX)

set(CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}" ${CMAKE_MODULE_PATH})
find_package(Tensorflow)

add_executable(TensorflowPlayground main.cpp)
target_include_directories(TensorflowPlayground PRIVATE ${TENSORFLOW_INCLUDE_DIR})
target_link_libraries(TensorflowPlayground ${TENSORFLOW_LIBRARY})
