#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin:/usr/local/lib/python2.7/bin:/opt/homebrew/bin
curPath=`pwd`
echo "web pnpm build start"


echo $curPath

if [ -f $curPath/fastcdn-web/apps/web-naive/dist.zip ] or [ -f $curPath/fastcdn/public/dist.zip ];then
	web_md5=`md5sum $curPath/fastcdn-web/apps/web-naive/dist.zip | awk '{print $1}'`
	fastcdn_public_md5=`md5sum $curPath/fastcdn/public/dist.zip | awk '{print $1}'`

	if [ "$web_md5" == "$fastcdn_public_md5" ];then
		# rm -rf ${curPath}/fastcdn-web/apps/web-naive/dist.zip
		echo "web file no change!"
		exit 0
	fi
fi

cd fastcdn-web && pnpm build
cp -f ${curPath}/fastcdn-web/apps/web-naive/dist.zip ${curPath}/fastcdn/public/dist.zip

cd ${curPath}/fastcdn/public
unzip -o dist.zip -d ./ 

echo "rm -rf ${curPath}/fastcdn-web/apps/web-naive/dist.zip"
rm -rf ${curPath}/fastcdn-web/apps/web-naive/dist.zip

echo "web pnpm build end"
