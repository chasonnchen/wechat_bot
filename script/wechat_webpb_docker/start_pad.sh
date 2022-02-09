#!/bin/bash

export WECHATY_LOG="verbose"
export WECHATY_PUPPET="wechaty-puppet-padlocal"
export WECHATY_PUPPET_PADLOCAL_TOKEN="puppet_padlocal_xxx"
export WECHATY_TOKEN="xxx"
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER=true
# Set port for your puppet service: must be published accessible on the internet
export WECHATY_PUPPET_SERVER_PORT=xxx

nohup docker run \
--name wechaty_puppet_service_pad2 \
--rm \
-ti \
-e WECHATY_LOG \
-e WECHATY_PUPPET \
-e WECHATY_PUPPET_PADLOCAL_TOKEN \
-e WECHATY_PUPPET_SERVER_PORT \
-e WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER \
-e WECHATY_TOKEN \
-p "$WECHATY_PUPPET_SERVER_PORT:$WECHATY_PUPPET_SERVER_PORT" \
wechaty/wechaty:latest &
