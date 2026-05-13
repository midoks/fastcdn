#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin


# curl -fsSL  https://raw.githubusercontent.com/midoks/fastcdn/master/scripts/install_dev.sh | sh

# Linux 手动安装
# wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz
# sudo tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz
# sudo ln -s /usr/local/go/bin/* /usr/bin/

# systemctl status fastcdn

# 手动编译
# go build main.go -o fastcdn && fastcdn web 

if [ ! -d /usr/local/go ];then
	wget https://golang.google.cn/dl/go1.26.2.linux-amd64.tar.gz
	tar -xvf go1.26.2.linux-amd64.tar.gz
	mv go /usr/local/
fi


# Debug Now
export PATH=/usr/local/go:$PATH:/root/go/bin
export GOPATH=/root/go


TAGRT_DIR=/usr/local/fastcdn_dev
mkdir -p $TAGRT_DIR
cd $TAGRT_DIR

export GIT_COMMIT=$(git rev-parse HEAD)
export BUILD_TIME=$(date -u '+%Y-%m-%d %I:%M:%S %Z')


if [ ! -d $TAGRT_DIR/fastcdn ]; then
	git clone https://github.com/midoks/fastcdn
	cd $TAGRT_DIR/fastcdn
else
	cd $TAGRT_DIR/fastcdn
	git pull
fi

go mod tidy
go mod vendor

# cd /usr/local/fastcdn_dev/fastcdn && go build -o fastcdn main.go 
# cd /usr/local/fastcdn_dev/fastcdn && go build -o fastcdn main.go && ./fastcdn web
cd $TAGRT_DIR/fastcdn && go build -o fastcdn main.go 
systemctl daemon-reload


cd $TAGRT_DIR/fastcdn && ./fastcdn install
systemctl restart fastcdn

cd $TAGRT_DIR/fastcdn && ./fastcdn -v

systemctl status fastcdn



# systemctl status fastcdn
# journalctl -u mgo_web -f
# systemctl stop fastcdn
# systemctl restart fastcdn