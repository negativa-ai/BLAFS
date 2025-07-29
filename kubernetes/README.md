# Kubernetes setup

This defines the files needed to run `baffs` in Kubernetes.

## Images

It derives everything from Minikube docker images.
As this runs a full featured Ubuntu 22.04 LTS inside,
it is compatible with the BLAFS containers.


The actual Dockerfile can reuse a lot of the original Minikube image.
The compile task has been separated in an early build stage derived
from the very same image used by Minikube for compatibility reasons
(the Golang build container artefacts were not compatible and the 
`docker run ...` did not succeed on a shadowed container)

The Docker version inside is used, not need to install another version.

## How to build the Minikube container with baffs inside

Build the image with `docker build .  -t minikube-baffs` and run Minikube
with 

```bash
minikube start --driver=docker --cpus=8 --memory=16384  --base-image=minikube-baffs`. 
```

In a terminal run `podeventlistener.sh`. That is all!

## How it works

In the Minikube container there is a Systemd service `unbloatpipe.service`
installed responsible for the shadowing process.

- __pulling the image__
- waiting for the image
- __creating the shadow image with `baffs`__
- waiting the image to be
  - used 
  - and freed
- __debloating the image with `baffs`__ 

The usage of the image is started from the Kubernetes side creating a pod and deleting it later by patching the object owning the pod.

This is triggered by the `podeventlistener.sh` script. It listens to the event pipeline
for the creation of pod, finds the owning object, normally a deployment or a daemonset
and patches it appending the suffix `-baffs` to the image version.

This way an update is triggered which can only be successfull after the debloat
process in the unbloatpipe is done. In between while the new container
might not exist the old pod is continued.

It uses `writepipe.yaml` to create a job connecting to the Minikube container

## Usage

You now can create a deplopment using this image by

```bash
kubectl create deployment redis-shadow --image=redis:7.4.1
```

Test your application as much as possible. This is emulated by a waiting period
of 20s for demo purposes.

The object owning the pod is automatically patched and will replace the old image by
deleting the initial pod.

## Testing

The `test.sh` script starts several deployment to show how it works and
if the scripts are successfully recovering from restarting the services
in Minikube.

## Maturity

This is a proof of concept. __Do not use it in production or for anything important.__
The `baffs` script restarts the docker daemon for each call and the `unbloatpipe`
additionally kills the Kubelet process.

To be ready for production it needs to be parallized and `baffs` must avoid to
restart the daemon.

