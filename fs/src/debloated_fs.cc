#define FUSE_USE_VERSION 36

#include <filesystem>
#include <stddef.h>

#include <fuse.h>
#include <spdlog/cfg/env.h>
#include <spdlog/spdlog.h>

#include "debloated_fs.h"
#include "util.h"

#define BAFFS_FUSE_OPTION(t, p) \
  {t, offsetof(struct Options, p), 1}

// Initialize the global FUSE_OPTS with the given arguments.
// `args` must be a valid pointer to a `fuse_args` struct.
// Files viewed from the mount point are actually located in `lower_dir`.
// Once accessed, they will be moved to `real_dir`.
void fuse_opt_init(struct fuse_args *args)
{

  const struct fuse_opt option_spec[] = {
      BAFFS_FUSE_OPTION("--realdir=%s", real_dir), BAFFS_FUSE_OPTION("--lowerdir=%s", lower_dir),
      BAFFS_FUSE_OPTION("--optimize=%s", optimize), FUSE_OPT_END};

  if (fuse_opt_parse(args, &BAFFS_FUSE_OPTS, option_spec, NULL) == -1)
  {
    spdlog::error("Failed to parse options");
    return;
  }

  spdlog::debug("lowerdir={0}, realdir={1}, optimize={2}", BAFFS_FUSE_OPTS.lower_dir,
                BAFFS_FUSE_OPTS.real_dir, BAFFS_FUSE_OPTS.optimize);
}

static const char *redirect(const char *original_path)
{
  spdlog::debug("original path: {0}", original_path);
  const char *lower_path = to_target_path(original_path, BAFFS_FUSE_OPTS.lower_dir);
  const char *real_path = to_target_path(original_path, BAFFS_FUSE_OPTS.real_dir);

  if (path_exists(real_path))
  {
    // file already moved to real dir, return real path
    spdlog::debug("real path exists, return real path: {0}", real_path);
    return real_path;
  }

  if (path_exists(lower_path))
  {
    spdlog::debug("lower path exists: {0}", lower_path);
    struct stat file_stat;
    lstat(lower_path, &file_stat);
    int rc = 0;
    if (S_ISDIR(file_stat.st_mode))
    {
      spdlog::debug("lower file is dir, create a new one: {0}->{1}, mode: {2}",
                    lower_path, real_path, file_stat.st_mode);
      // Theoretically, what we should do is a dir without its content
      // mkdir doesn't use st_mode directly
      // https://stackoverflow.com/questions/39737609/why-cant-my-program-set-0777-mode-with-the-mkdir-system-call

      if ((rc = mkdir(real_path, file_stat.st_mode)) != 0)
      {
        spdlog::error("error when creating dir: {0}; error code: {1}", real_path, rc);
      }
      if ((rc = chmod(real_path, file_stat.st_mode)) != 0)
      {
        spdlog::error("error when changing mode of dir: {0}; error code: {1}", real_path, rc);
      }
    }
    else
    {
      spdlog::debug("lower file is a regular file, copy it: {0}->{1}",
                    lower_path, real_path);
      if ((rc = copy_file(lower_path, real_path)) != 0)
      {
        spdlog::error("error when copying file: {0}->{1}; error code: {2}", lower_path, real_path, rc);
      }
    }
    if ((rc = copy_file_times(&file_stat, real_path)) != 0)
    {
      spdlog::error("error when setting file times: {0}->{1}; error code: {2}", lower_path, real_path, rc);
    }
    return real_path;
  }
  else
  {
    spdlog::debug("full lower path NOT exist: {0}", lower_path);
    return lower_path;
  }
}

int baffs_getattr(const char *_path, struct stat *stbuf,
                  struct fuse_file_info *fi)
{

  return 0;
}
