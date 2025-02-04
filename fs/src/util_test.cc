#include <gtest/gtest.h>
#include "util.h"

TEST(ToTargetPath, HandleNormalInput)
{
    const char *path = "path";
    const char *target_dir = "/lower";
    const char *expected = "/lower/path";
    const char *result = to_target_path(path, target_dir);
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
