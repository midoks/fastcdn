#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/usr/local/lib/python2.7/bin:/opt/homebrew/bin
curPath=`pwd`


# fastcdn
mkdir -p dist
mkdir -p dist/fastcdn
mkdir -p dist/fastcdn/logs
mkdir -p dist/fastcdn/configs
mkdir -p dist/fastcdn/web

cd $curPath/fastcdn && go build -o fastcdn main.go
cp -rf $curPath/fastcdn/fastcdn $curPath/dist/fastcdn/fastcdn

mkdir -p $curPath/dist/fastcdn/fastcdn-api
mkdir -p $curPath/dist/fastcdn/fastcdn-api/configs
mkdir -p $curPath/dist/fastcdn/fastcdn-api/bin
mkdir -p $curPath/dist/fastcdn/fastcdn-api/deploy
mkdir -p $curPath/dist/fastcdn/fastcdn-api/logs
mkdir -p $curPath/dist/fastcdn/fastcdn-api/data

# fastcdn-node
cd $curPath && mkdir -p dist/fastcdn-node/bin
cd $curPath mkdir -p dist/fastcdn-node/configs



cd $curPath/dist/fastcdn && ./fastcdn web

