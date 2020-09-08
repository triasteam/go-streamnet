#!/bin/bash

if true
then
  hostip=$(echo $HOST_IP|awk '{print $3}')
  echo $hostip
  sed -i "6s/localhost/${hostip}/g" config.yml
fi
go run main.go
