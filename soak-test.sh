#!/bin/bash
set -ex

# 2021-11-22T22:39:21.245939463Z
DATE=`date +"%Y-%m-%d"`
DATE=""$DATE"T22:39:21.245939463Z"
DATE2=`date +"%Y-%m-%d"`
DATE2=""$DATE2"T22:40:51.245939463Z"
for i in {0..10}
do
  cat example_web_hook_alert_message.json |\
  jq '.alerts[0].generatorURL = $newVal' --arg newVal "job-$i" |\
  jq '.alerts[0].startsAt = $dateVal' --arg dateVal "$DATE" |\
  jq '.alerts[0].endsAt = $endsat' --arg endsat "$DATE2" |\
  curl -X POST http://localhost:8801 -H "Content-Type: application/json" --data-binary @-
done

