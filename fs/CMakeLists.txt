cmake_minimum_required(VERSION 3.14)
project(BAFFS VERSION 0.1 LANGUAGES CXX)

# Enable C++11 (or later)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)


find_package(PkgConfig REQUIRED)
pkg_check_modules(FUSE3 REQUIRED fuse3)

include(FetchContent)
FetchContent_Declare(
  spdlog
  GIT_REPOSITORY https://github.com/gabime/spdlog.git
  GIT_TAG        v1.15.0
)
FetchContent_MakeAvailable(spdlog)



# Create the executable from our source file.
add_executable(debloated_fs src/debloated_fs.cc)
target_include_directories(debloated_fs PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include/)

# Add the FUSE3 include directory to the target
target_include_directories(debloated_fs PRIVATE ${FUSE3_INCLUDE_DIRS})
target_link_libraries(debloated_fs PRIVATE ${FUSE3_LIBRARIES})

# Add the spdlog include directory to the target
target_link_libraries(debloated_fs PRIVATE spdlog::spdlog)
