#!/bin/sh

export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

serverAdr="114.114.114.114"

ping -c 1 $serverAdr > /dev/null 2>&1
while [ $? -ne 0 ]; do
  echo -e "\e[1A\e[K $(date): Connecting - ${serverAdr}"
  sleep 1
  ping -c 1 $serverAdr > /dev/null 2>&1
done

echo "$(date): Connected - ${serverAdr}";

cd /home/aaron/dev/project/zv
git fetch --all && git reset --hard origin/main && git pull
go get all

# nohup go run cmd/vm/main.go -t vm > zv.log 2>&1
go run cmd/vm/main.go -t vm
