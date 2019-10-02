#!/bin/sh

npm install
npm run build

rm flashrom-prod
if [ "$1" = "rpi" ]; then
	export GOARM=6
	export GOARCH=arm
fi
go build

rice -i . append --exec flashrom-prod
