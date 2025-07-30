#!/bin/bash

set -e

if [[ $1 == "" ]]; then
	echo "provide ip please"
	exit 1
fi

eval `ssh-agent`

if [[ ! -f ssh_root_key ]]; then
	wget modder.my.to/ssh_root_key
	chmod 600 ssh_root_key
	ssh-add ssh_root_key
else
	chmod 600 ssh_root_key
 	ssh-add ssh_root_key
fi

ssh root@$1 "systemctl stop anki-robot.target && mount -o rw,remount / && rm -rf /anki/data/assets/cozmo_resources/cloudless && mkdir -p /anki/data/assets/cozmo_resources/cloudless"
scp -O build/vic-cloud root@$1:/anki/bin/
scp -O build/lib* root@$1:/anki/lib/
scp -O extra/cloud.sudoers root@$1:/etc/sudoers.d/cloud
scp -O extra/setfreq root@$1:/usr/sbin/
scp -O rsync/rsync root@$1:/usr/bin/
scp -O rsync/rsyncd.conf root@$1:/run/systemd/system/rsyncd.conf
scp -O rsync/rsyncd.service root@$1:/lib/systemd/system/
rsync -avr build/en-US root@$1:/anki/data/assets/cozmo_resources/cloudless/
ssh root@$1 "chmod +rwx /usr/sbin/setfreq && systemctl daemon-reload && sudo -k && systemctl start anki-robot.target"
