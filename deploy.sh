#!/bin/bash

set -e

if [[ $1 == "" ]]; then
	echo "provide ip please"
	exit 1
fi

ssh root@$1 "systemctl stop anki-robot.target && mount -o rw,remount / && rm -rf /anki/data/assets/cozmo_resources/cloudless && mkdir -p /anki/data/assets/cozmo_resources/cloudless"
scp build/vic-cloud root@$1:/anki/bin/
scp build/lib* root@$1:/anki/lib/
rsync -avr build/en-US root@$1:/anki/data/assets/cozmo_resources/cloudless/
