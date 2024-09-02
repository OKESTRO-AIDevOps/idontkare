# idontkare

Now I don't care about multi-cluster Kubernetes project management, because this one does 





## Get started

### server start

```shell

# start db in container

cd db && sudo docker compose up --build -d


# build client, server

make build


# build agent 

make build-agent


# start server 

cd src/server && ./server.out

```

### client example

```shell

cd src/client


# user set

./client.out user set --name test --pass test

# or use file for user set

./client.out user set --from-file $FILE_PATH

# file content should look like this

name: something
pass: somethingsecret

# cluster set, this will write private key pem to stdout

./client.out cluster set --username test --name testcluster

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

# modify config.yaml accordingly

# save private key content 

# connect agent

./agent.out

```


## Reference

This project has history.

See [nkia](https://github.com/OKESTRO-AIDevOps/nkia)

Whenever I think about whether this project is meaningful at all, I play [this song](https://www.youtube.com/watch?v=GVKRqIDS3WY) 