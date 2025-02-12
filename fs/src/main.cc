// MIT License

// Copyright (c) [2025] [jzh18]

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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