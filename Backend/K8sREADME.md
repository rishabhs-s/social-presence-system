kubectl scale deployment app-deployment --replicas=0 scale the deployment or it restarts

kubectl get svc- gets services

kubectl get pods- gets the pods

kubectl apply -f deployment.yaml - after change in yaml file to update change

kubectl delete pod backendapp1-688d4cf868-jhtnl - delete particular pod

eval $(minikube docker-env)- changes to minikube env-> build image->since it has its own docker env and cant access docker image on system

eval $(minikube docker-env -u) -- unsets the image

IMP- service name should be same for connection- In my case postgres-db is the service name which should be called in backenddeployment. 