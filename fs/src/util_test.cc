// MIT License

// Copyright (c) [2025] [jzh18]

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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
