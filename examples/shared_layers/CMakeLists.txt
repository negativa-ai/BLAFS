cmake_minimum_required(VERSION 3.10)

project(file_reader)

option(STATIC "build static library" ON)

add_executable(file_reader main.cc)
# build the executable staticly
if(STATIC)
	message("build static library")
	target_link_libraries(file_reader PUBLIC "-static")
endif()