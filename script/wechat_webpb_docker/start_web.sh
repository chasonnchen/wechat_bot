#!/bin/bash
export WECHATY_LOG="verbose"
export WECHATY_PUPPET="wechaty-puppet-wechat"
export WECHATY_PUPPET_SERVER_PORT="30005"
export WECHATY_TOKEN="2fdb00a5-5c31-4018-84ac-c64e5f995057"
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER=true

nohup docker run -ti \
--name wechaty_puppet_service_token_gateway \
--rm \
-e WECHATY_LOG \
-e WECHATY_PUPPET \
-e WECHATY_PUPPET_SERVER_PORT \
-e WECHATY_TOKEN \
-e WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER \
-p "$WECHATY_PUPPET_SERVER_PORT:$WECHATY_PUPPET_SERVER_PORT" \
wechaty/wechaty:latest &
