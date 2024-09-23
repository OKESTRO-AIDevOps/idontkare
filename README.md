# idontkare

Now I don't care about multi-cluster Kubernetes project management, because this one does 





## Get started

### requirements

- linux
- go
- docker
- make
- kind (kubernets in docker, for development)

### server start

```shell

# start db in container

cd db && sudo docker compose up --build -d


# build client, server

make build SERVER_HOST=localhost # or whatever the domain name of yours


# build agent 

make build-agent


# start server 

cd src/server && ./server.out

```

### mock cluster for development

```shell

# kubernetes cluster in docker


cd hack/cluster && sudo ./kindcluster.sh $USER

# check if cluster is up and running

kubectl get nodes

# to destroy

sudo kind delete cluster --name kindcluster

```


### client example

```shell

cd src/client


# user set

./client.out user set --name sampleusername --pass sampleuserpass

# or use file for user set

./client.out user set --from-file $FILE_PATH

# file content should look like this

name: something
pass: somethingsecret

# cluster set, this will write private key pem to stdout

./client.out cluster set --username sampleusername --name samplecluster

# project set 

./client.out project set --from-file $FILE_PATH

# update ci option

./client.out project ci option set --username sampleusername --name sampleproject --path ./sample/cioption.yaml

# update cd option

./client.out project cd option set --username sampleusername --name sampleproject --path ./sample/cdoption.yaml

```

### agent example

```shell

cd src/agent


# save private key content 

# modify config.yaml accordingly

# connect agent

./agent.out

```




## Reference

This project has history.

See [nkia](https://github.com/OKESTRO-AIDevOps/nkia)

Whenever I think about whether this project is meaningful at all, I play [this song](https://www.youtube.com/watch?v=GVKRqIDS3WY) 