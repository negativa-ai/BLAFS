#include <utime.h>
#include <sys/stat.h>

#include <filesystem>
#include <string>

#include "util.h"

// concat path1 and path2 with a '/' in between
std::string concat_path(const char *path1, const char *path2)
{
    std::filesystem::path target{path1};
    target += "/";
    target += path2;
    return target.lexically_normal().string();
}

std::string concat_path(const char **paths, uint32_t size)
{
    if (size == 0)
    {
        return "";
    }

    std::filesystem::path target{paths[0]};
    for (uint32_t i = 1; i < size; i++)
    {
        target += "/";
        target += paths[i];
    }
    return target.lexically_normal().string();
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