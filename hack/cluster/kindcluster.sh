#!/bin/bash

if [ -z "$1" ]
then

    echo "needs \$USERNAME"

    exit 1

fi

USERNAME="$1"

sudo kind create cluster --name kindcluster --config ./kindcluster.yaml --image=kindest/node:v1.27.2

sudo /bin/cp -Rf /root/.kube/config /home/$USERNAME/.kube/config 

sudo chown $USERNAME:$USERNAME /home/$USERNAME/.kube/config 