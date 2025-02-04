#define FUSE_USE_VERSION 36

#include <filesystem>
#include <stddef.h>
#include <string.h>

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

std::string redirect(const char *original_path)
{
  spdlog::debug("original path: {0}", original_path);
  std::string lower_path_string = to_target_path(original_path, BAFFS_FUSE_OPTS.lower_dir);
  const char *lower_path = lower_path_string.c_str();
  std::string real_path_string = to_target_path(original_path, BAFFS_FUSE_OPTS.real_dir);
  const char *real_path = real_path_string.c_str();

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

void *baffs_init(struct fuse_conn_info *conn, struct fuse_config *cfg)
{
  cfg->use_ino = 1;

  /* Pick up changes from lower filesystem right away. This is
     also necessary for better hardlink support. When the kernel
     calls the unlink() handler, it does not know the inode of
     the to-be-removed entry and can therefore not invalidate
     the cache of the associated inode - resulting in an
     incorrect st_nlink value being reported for any remaining
     hardlinks to this inode. */
  cfg->entry_timeout = 0;
  cfg->attr_timeout = 0;
  cfg->negative_timeout = 0;

  return NULL;
}

int baffs_getxattr(const char *_path, const char *name, char *value, size_t size)
{
  spdlog::debug("getxattr callback, not supported");
  return -ENOTSUP;
}

int baffs_listxattr(const char *_path, char *list, size_t size)
{
  spdlog::debug("listxattr callback, not supported");
  return -ENOTSUP;
}

int baffs_getattr(const char *_path, struct stat *stbuf,
                  struct fuse_file_info *fi)
{
  std::string redirected_path_string = redirect(_path);
  const char *redirected_path = redirected_path_string.c_str();

  if (fi)
  {
    spdlog::debug("retrieve from fi");
    return fstat(fi->fh, stbuf);
  }
  else
  {
    spdlog::debug("retrieve from path: {0}", redirected_path);
    return lstat(redirected_path, stbuf);
  }
}
