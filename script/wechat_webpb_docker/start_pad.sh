#!/bin/bash

export WECHATY_LOG="verbose"
export WECHATY_PUPPET="wechaty-puppet-padlocal"
export WECHATY_PUPPET_PADLOCAL_TOKEN="puppet_padlocal_7cbebc1ab76f41a58cfd25c0ff3eaf4b"
export WECHATY_TOKEN="2fdb00a5-5c31-4018-84ac-c64e5f995057"
export WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_SERVER=true
# Set port for your puppet service: must be published accessible on the internet
export WECHATY_PUPPET_SERVER_PORT=30009

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
