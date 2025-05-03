# 測試環境

補充中

``` docker[compile]
sudo docker run -v /home/chris/git:/app -v /home/chris/git/wander/env_test/dist:/dist -w /app --rm -ti --name build_go golang:1.24-bullseye bash

rm /dist/* -fr && cd /app/wander/env_test && sh build.sh
```

``` sh
cd /home/chris/git/wander/env_test

sudo docker-compose -f docker-compose-init.yaml up -d

sudo docker-compose -f docker-compose-server.yaml up -d
```
