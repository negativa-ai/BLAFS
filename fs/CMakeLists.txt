cmake_minimum_required(VERSION 3.14)
project(BAFFS VERSION 0.1 LANGUAGES CXX)

# Enable C++11 (or later)
set(CMAKE_CXX_STANDARD 17)
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

FetchContent_Declare(
  googletest
  GIT_REPOSITORY https://github.com/google/googletest.git
  GIT_TAG        v1.15.2
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

# Create the executable from our source file.
add_executable(debloated_fs src/main.cc src/debloated_fs.cc)
target_include_directories(debloated_fs PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include/)

# Add the FUSE3 include directory to the target
target_include_directories(debloated_fs PRIVATE ${FUSE3_INCLUDE_DIRS})
target_link_libraries(debloated_fs PRIVATE ${FUSE3_LIBRARIES})

# Add the spdlog include directory to the target
target_link_libraries(debloated_fs PRIVATE spdlog::spdlog)

# google test

# Create the executable from our source file.
add_executable(all_test src/util_test.cc src/util.cc)
target_include_directories(all_test PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include/)

# Add the FUSE3 include directory to the target
target_include_directories(all_test PRIVATE ${FUSE3_INCLUDE_DIRS})
target_link_libraries(all_test PRIVATE ${FUSE3_LIBRARIES})

# Add the spdlog include directory to the target
target_link_libraries(all_test PRIVATE spdlog::spdlog)

# Add the gtest include directory to the target
target_link_libraries(all_test PRIVATE gtest_main)

include(GoogleTest)
gtest_discover_tests(all_test)
