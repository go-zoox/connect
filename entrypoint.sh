#!/bin/sh

export MODE=${MODE:-"production"}
export PORT=${PORT:-"8080"}
export CONFIG=${CONFIG:-"/conf/config.yml"}

function main() {
  connect serve
}

main
