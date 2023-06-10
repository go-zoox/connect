#!/bin/sh

function main() {
  if [ "$MODE" = "doreamon" ]; then
    connect doreamon
    if [ "$?" != "0" ]; then
      return $?
    fi

    return
  fi

  connect serve \
    -p ${PORT:-8080} \
    -c ${DIR:-"/conf/config.yml"}
}

main
