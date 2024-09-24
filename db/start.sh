#!/bin/bash


if [ -z "$1" ]
then

    echo "needs \$USERNAME"

    exit 1

fi

USERNAME="$1"

sudo docker compose up --build -d

sudo chown -R "$USERNAME:$USERNAME" data