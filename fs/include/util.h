
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
#ifndef BAFFS_FS_UTIL_H_
#define BAFFS_FS_UTIL_H_

#include <string>
// Concats path1 and path2 with a '/' in between
std::string concat_path(const char *path1, const char *path2);

// Concats paths with a '/' in between.
// The size is the number of paths in the array.
std::string concat_path(const char **paths, uint32_t size);

// Returns true if the path exists.
bool path_exists(const char *path);

// Copies file from `src_file` to `dst_file`
int copy_file(const char *from, const char *to);

// Copies the timestamp info of `src_file` to `dst_path`
int copy_file_times(struct stat *src_file, const char *dst_path);

#endif // BAFFS_FS_UTIL_H_
