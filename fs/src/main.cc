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
        .mknod = baffs_mknod,
        .mkdir = baffs_mkdir,
        .unlink = baffs_unlink,
        .rmdir = baffs_rmdir,
        .symlink = baffs_symlink,
        .rename = baffs_rename,
        .link = baffs_link,
        .chmod = baffs_chmod,
        .chown = baffs_chown,
        .truncate = baffs_truncate,
        .open = baffs_open,
        .read = baffs_read,
        .write = baffs_write,
        .statfs = baffs_statfs,
        .flush = baffs_flush,
        .release = baffs_release,
        .fsync = baffs_fsync,
        .getxattr = baffs_getxattr,
        .listxattr = baffs_listxattr,
        .opendir = baffs_opendir,
        .readdir = baffs_readdir,
        .releasedir = baffs_releasedir,
        .init = baffs_init,
        .access = baffs_access,
        .create = baffs_create,
        .ioctl = baffs_ioctl,
        .lseek = baffs_lseek,

    };

    int ret = fuse_main(args.argc, args.argv, &baffs_oper, nullptr);
    fuse_opt_free_args(&args);
    return ret;
}