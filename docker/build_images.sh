#!/bin/sh

set -x

work_space=".."
version=$1

# 是否指定版本号
if [ -z "$version" ]; then
  version="latest"
fi

cp -r $work_space/bin $work_space/docker/bin

tag=swr.cn-east-3.myhuaweicloud.com/badminton/feishu_doc_blocks_plugin_server:$version
docker build --platform linux/amd64 -t $tag -f Dockerfile .

rm -rf $work_space/docker/bin
