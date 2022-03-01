@echo on

:start

ping /n 3 114.114.114.114 | findstr "TTL=" && goto next || goto start

:next
cmd /k "cd /d C:\Users\admin\dev\project\zv && git fetch --all && git reset --hard origin/main && git pull && go env -w GOPROXY=https://goproxy.io && go get all && go run cmd\vm\main.go -t vm -s http://192.168.0.232:8455"

pause
