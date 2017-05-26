#!/bin/bash
while :
do
	mosquitto_pub -t 'orders' -m 'helloWorld'
	mosquitto_pub -t 'requests' -m 'helloWorld'
	mosquitto_pub -t 'bargains' -m 'helloWorld'
	sleep 1
done
