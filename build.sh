#!/bin/sh
cd ./web
npm install
npm run build
cd ../
rm -rf ./md/web
cp ./web/dist -r ./md/web
cd ./md
go build
echo '编译完成'
