#define FUSE_USE_VERSION 36

#include <fuse.h>
#include <spdlog/cfg/env.h>
#include <spdlog/spdlog.h>

#include "debloated_fs.h"

int main(int argc, char **argv)
{
    spdlog::cfg::load_env_levels();

    struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

    fuse_opt_init(&args);

    const struct fuse_operations baffs_oper = {
        .getattr = baffs_getattr,
        .readlink = baffs_readlink,
        .getxattr = baffs_getxattr,
        .listxattr = baffs_listxattr,
        .opendir = baffs_opendir,
        .readdir = baffs_readdir,
        .init = baffs_init,
        .access = baffs_access,

    };

    int ret = fuse_main(args.argc, args.argv, &baffs_oper, nullptr);
    fuse_opt_free_args(&args);
    return ret;
}