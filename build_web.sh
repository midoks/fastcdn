#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/usr/local/lib/python2.7/bin:/opt/homebrew/bin
curPath=`pwd`
echo "web pnpm build start"


web_md5=`md5sum fastcdn-web/apps/web-naive/dist.zip | awk '{print $1}'`
fastcdn_public_md5=`md5sum fastcdn/public/dist.zip | awk '{print $1}'`

if [ "$web_md5" == "$fastcdn_public_md5" ];then
	echo "web file no change!"
	exit 0
fi

cd fastcdn-web && pnpm build
cp -f apps/web-naive/dist.zip ../fastcdn/public/dist.zip
cd ../fastcdn/public
unzip -o dist.zip -d ./ 


echo "web pnpm build end"
