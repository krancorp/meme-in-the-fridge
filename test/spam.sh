#!/bin/bash
while :
do
	echo -n "Krombacher:2" >/dev/udp/$1/8080
done
