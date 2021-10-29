## Kubernetes custom controller 
Creation of services and Ingress controller is going to trigger automatically in case of a manual Deployment creation in a given namespace or accross all the namespaces .

## Usages

create a Kubernetes deployment
```$ kubectl create deployment nginx --image nginx```

watch for the service creation
```$ watch kubectl get svc```

forward the container port 
```$ kubectl port-forward svc/<name> 8080:80```