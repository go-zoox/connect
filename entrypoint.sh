#!/bin/sh

server \
  -p ${PORT:-8080} \
  -c ${DIR:-"/conf/config.yml"}
