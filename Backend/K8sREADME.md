kubectl scale deployment app-deployment --replicas=0 scale the deployment or it restarts

kubectl get svc- gets services

kubectl get pods- gets the pods

kubectl apply -f deployment.yaml - after change in yaml file to update change

kubectl delete pod backendapp1-688d4cf868-jhtnl - delete particular pod

eval $(minikube docker-env)- changes to minikube env-> build image->since it has its own docker env and cant access docker image on system

eval $(minikube docker-env -u) -- unsets the image

IMP- service name should be same for connection- In my case postgres-db is the service name which should be called in backenddeployment. 

minikube ip - gets the ip


To expose the port on which deployment is running

kubectl get deployments
NAME                  READY   UP-TO-DATE   AVAILABLE   AGE
app-deployment        0/0     0            0           36m
backendapp            1/1     1            1           54m
backendapp1           0/0     0            0           31m
postgres-db           1/1     1            1           24m
postgres-deployment   0/0     0            0           38m
rishabhsharma@Rishabhs-MacBook-Air Backend % minikube kubectl expose deployment backendapp -- --type=NodePort --port=8080
then port forward local to pod
kubectl port-forward backendapp-56d6f4964f-kc6gt   8080:8080


To run on minikube- just kubectl apply -f dep.yaml(all yaml)-> services and deployments are created. Then Expose the port. 