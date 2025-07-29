
kubectl create deployment redis-shadow-1 --image=redis:7.4.1
kubectl create deployment redis-shadow-2 --image=redis:7.4.2

kubectl create deployment nginx-shadow-23 --image=nginx:1.23.0
kubectl create deployment nginx-shadow-24 --image=nginx:1.24.0
kubectl create deployment nginx-shadow-25 --image=nginx:1.25.0

kubectl create deployment redis-shadow-3 --image=redis:7.4.3
kubectl create deployment redis-shadow-4 --image=redis:7.4.4
kubectl create deployment redis-shadow-5 --image=redis:7.4.5

