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
#ifndef BAFFS_FS_DEBLOATED_FS_H_
#define BAFFS_FS_DEBLOATED_FS_H_

#define FUSE_USE_VERSION 36

#include <dirent.h>

#include <fuse.h>

#ifdef linux
/* For pread()/pwrite()/utimensat() */
#define _XOPEN_SOURCE 700
#endif

// Command line options of the executable
struct Options
{
    const char *real_dir;  // files accessed will be moved to real_dir from lower_dir
    const char *lower_dir; // files exist in lower_dir originally
    const char *optimize;  // reserved for future optimization
};

// The global variable to store the command line options
static struct Options BAFFS_FUSE_OPTS{NULL, NULL, NULL}; // global variable to store options

// The structure to store the directory pointer and the offset
struct BaffsDirp
{
    DIR *dp;
    struct dirent *entry;
    off_t offset;
};

// Initialize the global `BAFFS_FUSE_OPTS`, must be called before mounting
void fuse_opt_init(struct fuse_args *args);

/**
 * All the functions below are the callbacks of FUSE operations
 */

void *baffs_init(struct fuse_conn_info *conn, struct fuse_config *cfg);

int baffs_getxattr(const char *_path, const char *name, char *value,
                   size_t size);

int baffs_listxattr(const char *_path, char *list, size_t size);

int baffs_getattr(const char *_path, struct stat *stbuf,
                  struct fuse_file_info *fi);

int baffs_access(const char *_path, int mask);

int baffs_readlink(const char *_path, char *buf, size_t size);

int baffs_opendir(const char *_path, struct fuse_file_info *fi);

int baffs_readdir(const char *_path, void *buf, fuse_fill_dir_t filler,
                  off_t offset, struct fuse_file_info *fi,
                  enum fuse_readdir_flags flags);

int baffs_releasedir(const char *_path, struct fuse_file_info *fi);

int baffs_mknod(const char *_path, mode_t mode, dev_t rdev);

int baffs_mkdir(const char *_path, mode_t mode);

int baffs_unlink(const char *_path);

int baffs_rmdir(const char *_path);

int baffs_symlink(const char *_from, const char *_to);

int baffs_rename(const char *_from, const char *_to, unsigned int flags);

int baffs_link(const char *_from, const char *_to);

int baffs_chmod(const char *_path, mode_t mode,
                struct fuse_file_info *fi);

int baffs_chown(const char *_path, uid_t uid, gid_t gid,
                struct fuse_file_info *fi);

int baffs_truncate(const char *_path, off_t size,
                   struct fuse_file_info *fi);

int baffs_create(const char *_path, mode_t mode,
                 struct fuse_file_info *fi);

int baffs_open(const char *_path, struct fuse_file_info *fi);

int baffs_read(const char *_path, char *buf, size_t size, off_t offset,
               struct fuse_file_info *fi);

int baffs_write(const char *_path, const char *buf, size_t size,
                off_t offset, struct fuse_file_info *fi);

int baffs_statfs(const char *_path, struct statvfs *stbuf);

int baffs_release(const char *_path, struct fuse_file_info *fi);

int baffs_fsync(const char *_path, int isdatasync,
                struct fuse_file_info *fi);
int baffs_flush(const char *_path, struct fuse_file_info *fi);
off_t baffs_lseek(const char *_path, off_t off, int whence,
                  struct fuse_file_info *fi);

int baffs_ioctl(const char *_path, unsigned int cmd, void *arg,
                struct fuse_file_info *fi, unsigned int flags,
                void *data);

#endif // BAFFS_FS_BAFFS_H_
