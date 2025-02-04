
#ifndef BAFFS_FS_UTIL_H_
#define BAFFS_FS_UTIL_H_

#include <string>

std::string to_target_path(const char *path, const char *target_dir);
bool path_exists(const char *path);
int copy_file(const char *from, const char *to);

// copy file times from src_file to dst_path
int copy_file_times(struct stat *src_file, const char *dst_path);

#endif // BAFFS_FS_UTIL_H_
