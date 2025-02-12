# BAFFS
![example workflow](https://github.com/jzh18/BAFFS/actions/workflows/main.yml/badge.svg)
Shrink you container size up to 95% with BAFFS!

## Introduction

BAFFS is a bloat-aware filesystem for container debloating.
The design principles of BAFFS are effective, efficient, and easy to use.


## Installation

The easiest way to install BAFFS is to use Docker image:
```
docker pull justinzhf/baffs:latest
```

## Quick Start

1. Pull the BAFFS image:
    ```
    docker pull justinzhf/baffs:latest
    ```
2. Run the BAFFS container with privileged mode, so that we can run Docker in Docker:
    ```
    docker run  -d --name baffs --privileged=true  -v /tmp/docker:/var/lib/docker justinzhf/baffs:latest
    ```
    Note that we mount the host's `/tmp/docker` to the container's `/var/lib/docker`. 
    In this way, all the images pulled insied the container `baffs` will be stored in the host's `/tmp/docker`.
    You can change the path to any other directory you like.
3. Enter shell of the container:
    ```
    docker exec -it baffs zsh
    ``` 
4. Inside the container, we pull a redis image. We will debloat this image later.
    ```
    docker pull redis:7.4.1
    ```
5. Now we can start to debloat the redis image. 
    ```
    baffs shadow --images=redis:7.4.1
    ```
    This will convert the filesystem of the redis image to BAFFS filesystem.
6. Now we run the redis container with profiling workload. For example, we simply start the redis server:
    ```
    docker run -it --name redis redis:7.4.1
    ```
    After the redis server is started, use `Ctrl+C` to stop the redis server.
7. At this step, BAFFS has detected all the files needed by the redis server. We can now debloat the redis image:
    ```
    baffs debloat --images=redis:7.4.1
    ```
    This will generate a new redis image named `redis:7.4.1-baffs`, which has a much smaller size.

8. Now let's compare the size of the redis image before and after debloating:
    ```
    docker images | grep redis
    ```
    Here is an example output:
    ```
    redis        7.4.1-baffs   d43e8b090126   4 months ago   28.8MB
    redis        7.4.1         2724e40d4303   4 months ago   117MB
    ```
9. Finally, let's check whether the debloated image can still run the redis server:
    ```
    docker run -it --name redis redis:7.4.1-baffs
    ```
    If the redis server can be started, then the debloating is successful!

## Advanced Usage

### Debloat Multiple Images at Once
To Be Added

### Debloat Certain Layers of an Image
To Be Added

## Citation
Please cite our paper if you use BAFFS in your research:
```
@misc{zhang2023blafsbloatawarefile,
      title={BLAFS: A Bloat Aware File System}, 
      author={Huaifeng Zhang and Mohannad Alhanahnah and Ahmed Ali-Eldin},
      year={2023},
      eprint={2305.04641},
      archivePrefix={arXiv},
      primaryClass={cs.SE},
      url={https://arxiv.org/abs/2305.04641}, 
}
```