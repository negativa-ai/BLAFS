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

## Howto use

Build the image with `docker build .  -t minikube-baffs` and run Minikube
with `minikube start --driver=docker --base-image=minikube-baffs`. 

In the Minikube container trigger download of the image and the initial
shadowin  by using `writepipe redis:7.4.1`

You now can create a deplopment using this image by

```bash
kubectl create deployment redis-shadow --image=redis:7.4.1
```

Test your application as much as possible.

Now stop the deployment by `kubect delete deployment redis-shadow` 
and retrigger the debloat by using `writepipe redis:7.4.1`
in Minikube again.

Restart the deployment, it is using the debloated container image now!
