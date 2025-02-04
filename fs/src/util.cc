#include <utime.h>
#include <sys/stat.h>

#include <filesystem>
#include <string>

#include "util.h"

// path must be an absolute path
std::string to_target_path(const char *path, const char *target_dir)
{
    std::filesystem::path target{target_dir};
    target += path;
    return target.string();
}

bool path_exists(const char *path)
{
    struct stat stbuf;
    return lstat(path, &stbuf) == 0;
}

int copy_file(const char *from, const char *to)
{
    std::string cmd{"cp -rP \"" + std::string(from) + "\" \"" + std::string(to) + "\""};
    return system(cmd.c_str());
}

int copy_file_times(struct stat *src_file, const char *dst_path)
{
    struct utimbuf new_times;
    new_times.actime = src_file->st_atime;
    new_times.modtime = src_file->st_mtime;
    return utime(dst_path, &new_times);
}