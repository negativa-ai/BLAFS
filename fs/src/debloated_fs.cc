#define FUSE_USE_VERSION 36

#include <stddef.h>

#include <fuse.h>
#include <spdlog/cfg/env.h>
#include <spdlog/spdlog.h>

#include "baffs.h"

#define BAFFS_FUSE_OPTION(t, p) \
  {t, offsetof(struct Options, p), 1}

int main(int argc, char **argv)
{
  spdlog::cfg::load_env_levels();

  struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

  const struct fuse_opt option_spec[] = {
      BAFFS_FUSE_OPTION("--realdir=%s", real_dir), BAFFS_FUSE_OPTION("--lowerdir=%s", lower_dir),
      BAFFS_FUSE_OPTION("--optimize=%s", optimize), FUSE_OPT_END};

  struct Options opts;

  if (fuse_opt_parse(&args, &opts, option_spec, NULL) == -1)
  {
    spdlog::error("Failed to parse options");
    return 1;
  }

  return 0;
}
