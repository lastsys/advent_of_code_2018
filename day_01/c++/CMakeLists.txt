cmake_minimum_required(VERSION 3.11)
project(day_01)

file(GLOB_RECURSE main_sources src/*.cc src/*.h)

add_executable(day_01 ${main_sources})
