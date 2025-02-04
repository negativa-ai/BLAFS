#define FUSE_USE_VERSION 36

#include <filesystem>
#include <stddef.h>

#include <fuse.h>
#include <spdlog/cfg/env.h>
#include <spdlog/spdlog.h>

#include "debloated_fs.h"


#define BAFFS_FUSE_OPTION(t, p) \
  {t, offsetof(struct Options, p), 1}

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

int baffs_getattr(const char *_path, struct stat *stbuf,
                  struct fuse_file_info *fi)
{

  return 0;
}
