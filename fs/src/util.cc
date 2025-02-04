#include <filesystem>

#include "util.h"

const char *to_target_path(const char *path, const char *target_dir)
{
    std::filesystem::path lower{target_dir};
    lower /= path;
    return lower.c_str();
}
