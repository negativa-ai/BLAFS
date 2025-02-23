# BLAFS

![example workflow](https://github.com/jzh18/BAFFS/actions/workflows/main.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Shrink you container size up to 95%.

## Introduction

BLAFS is a bloat-aware filesystem for container debloating.
The design principles of BLAFS are effective, efficient, and easy to use.
It detects the files used by the container, and then debloats the container by removing the unused files.
The debloated containers are still functional and can run the same workload as the original containers, but with a much smaller size and faster deployment.

Check the paper for more details: [The Cure is in the Cause: A Filesystem for Container Debloating](https://arxiv.org/abs/2305.04641).

## Installation

The easiest way to install BAFFS is to use the Docker image:
```
docker pull justinzhf/baffs:latest
```

## Quick Start

1. Pull the BLAFS image:
    ```
    docker pull justinzhf/baffs:latest
    ```
2. Run the BLAFS container with privileged mode, so that we can run Docker in Docker:
    ```
    docker run  -d --name baffs --privileged=true  -v /tmp/docker:/var/lib/docker justinzhf/baffs:latest
    ```
    Note that we mount the host's `/tmp/docker` to the container's `/var/lib/docker`. 
    In this way, all the images pulled inside the container `baffs` will be stored in the host's `/tmp/docker`.
    You can change the path to any other directory you like.
3. Enter shell of the container:
    ```
    docker exec -it baffs bash
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
    docker run -it --rm redis:7.4.1
    ```
    After the redis server is started, use `Ctrl+C` to stop the redis server.
7. At this step, BLAFS has detected all the files needed by the redis server. We can now debloat the redis image:
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
    docker run -it --rm redis:7.4.1-baffs
    ```
    If the redis server can be started, then the debloating is successful!

## Script BLAFS 


BLAFS debloats a container in three steps:  

1. **Convert** – Converts the container into the BLAFS filesystem.  
2. **Profiling** – Runs profiling workloads to track file usage.  
3. **Debloating** – Retains only the files used during profiling, removing everything else.  

This script provides an example of using BLAFS to debloat a Redis container.  
🔗 [Example Script](https://github.com/negativa-ai/BLAFS/blob/main/tests/test.sh)



## Advanced Usage

BLAFS has three working modes: no-sharing, sharing, and serverless. 
Please refer to the paper for more details.

### Set Logging Level
Set logging level for `baffs`:
```
LOG_LEVEL=debug|info|warning|error baffs ...
```
Set logging level for `debloated_fs`:
```
SPDLOG_LEVEL=debug|info|warning|error baffs ...
```

### Debloat Multiple Images at Once
If two images share some common layers, we can debloat them together.
And the debloated images will share the same layers.


```
baffs shadow --images=img1,img2 # shadow multiple images, the two images should share the some layers initially
# run the profiling workload for img1 and img2
baffs debloat --images=img1,img2 # this will debloat both img1 and img2, with shared layers
```

### Debloat Certain Layers of an Image
Serverless containers are usually built on top of a base image.
We can debloat the only the unique layers of the serverless container while keeping the base image untouched.

```
baffs shadow --images=img1
# run the profiling workload for img1
baffs debloat --images=img1 --top=3 # debloat img1 with top 3 layers
```

## Citation
Please cite our paper if you use BLAFS in your research:
```
@misc{zhang2025curecausefilesystemcontainer,
      title={The Cure is in the Cause: A Filesystem for Container Debloating}, 
      author={Huaifeng Zhang and Philipp Leitner and Mohannad Alhanahnah and Ahmed Ali-Eldin},
      year={2025},
      eprint={2305.04641},
      archivePrefix={arXiv},
      primaryClass={cs.SE},
      url={https://arxiv.org/abs/2305.04641}, 
}
```
