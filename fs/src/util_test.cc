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
