#!/bin/sh

export MODE=${MODE:-"production"}
export PORT=${PORT:-"8080"}

function main() {
  connect doreamon
}

main
