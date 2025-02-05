#include <gtest/gtest.h>
#include "util.h"

TEST(ConcatPath, HandleNormalInput)
{
    const char *path1 = "/lower";
    const char *path2 = "/path";
    const char *expected = "/lower/path";
    const char *result = concat_path(path1, path2).c_str();
    ASSERT_STREQ(expected, result);
}

TEST(ConcatPath, HandleRootDir)
{
    const char *path1 = "/lower";
    const char *path2 = "/";
    const char *expected = "/lower/";
    const char *result = concat_path(path1, path2).c_str();
    ASSERT_STREQ(expected, result);
}

TEST(ConcatPath, HandleMultiPaths)
{
    const char *paths[3] = {"/lower", "/", "path"};
    const char *expected = "/lower/path";
    const char *result = concat_path(paths, 3).c_str();
    ASSERT_STREQ(expected, result);
}

TEST(PathExists, HandlePathNotExists)
{
    const char *path = "path";
    ASSERT_FALSE(path_exists(path));
}

TEST(PathExists, HandlePathExists)
{
    const char *path = "/usr";
    ASSERT_TRUE(path_exists(path));
}
