# idontkare

Now I don't care about multi-cluster Kubernetes project management, because this one does 





## Get started

### requirements

- linux
- go
- docker
- make
- kind (kubernets in docker, for development)

### 1. server start

```shell

# start db in container

cd db && sudo docker compose up --build -d

# close db 

cd db && sudo docker compose down


# build client, server

make build SERVER_HOST=localhost # or whatever the domain name of yours


# build agent 

make build-agent


# start server 

cd src/server && ./server.out

```

### 1-a. mock cluster for development

```shell

# kubernetes cluster in docker


cd hack/cluster && sudo ./kindcluster.sh $USER

# check if cluster is up and running

kubectl get nodes

# to destroy

sudo kind delete cluster --name kindcluster

```


### 2. set up cluster and project

```shell

cd src/client


# user set

./client.out user set --name sampleusername --pass sampleuserpass

# or use file for user set

./client.out user set --from-file $FILE_PATH

# for example

./client.out user set --from-file ./sample/userset.yaml


# cluster set, this will write private key pem to stdout

./client.out cluster set --username sampleusername --name samplecluster

# project set 

./client.out project set --from-file ./sample/projectset.yaml


# the above won't work, obviously!
# to make it work
# replace content properly 



```

### 3. connect cluster

```shell

cd src/agent


# save private key content 

# modify config.yaml accordingly

# connect agent

./agent.out

```

### 4. test build and deployment


```shell

cd src/client

# update ci option

./client.out project ci option set --username sampleusername --name sampleproject --path ./sample/cioption.yaml

# update cd option

./client.out project cd option set --username sampleusername --name sampleproject --path ./sample/cdoption.yaml


```

### 5. more client requests 


```shell

# get all ci history for project

./client.out project ci history get all --username sampleusername --project sampleproject

# get all cd history for project

./client.out project cd history get all --username sampleusername --project sampleproject

# get lifecycle report for project

./client.out lifecycle report get latest --username sampleusername --project sampleproject

```


## Reference

This project has history.

See [nkia](https://github.com/OKESTRO-AIDevOps/nkia)

Whenever I think about whether this project is meaningful at all, I play [this song](https://www.youtube.com/watch?v=GVKRqIDS3WY) 