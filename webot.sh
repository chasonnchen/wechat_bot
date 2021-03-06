#!/bin/bash
# 	GitHb: https://github.com/chasonnchen/wechat_bot
# 	Author: Chasonn Chen <185320860@qq.com>
cmd=$1
if [ -n ${cmd} ]; then
    echo "Start to "${cmd}" the project."
else
     echo "Cmd error, eg: sh cmd.sh restart"
     exit
fi



if [ "$cmd" == "build" ];then
    go build
    rm -rf runtime/wechat_bot
    mv wechat_bot runtime/
fi

if [ "$cmd" == "kill" ];then
    ps -ef | grep wechat_bot | grep -v grep | awk '{print $2}' | xargs kill -9
fi

if [ "$cmd" == "restart" ];then
    ps -ef | grep wechat_bot | grep -v grep | awk '{print $2}' | xargs kill -9
    cd runtime
    nohup ./wechat_bot &
fi
