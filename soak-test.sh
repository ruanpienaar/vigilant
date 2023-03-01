#!/bin/bash
set -ex

for i in {0..10}
do
  DATE=`date -v-${i}d +"%Y-%m-%d"`
  cat example_web_hook_alert_message.json |\
  jq '.alerts[0].generatorURL = $newVal' --arg newVal "$i" |\
  jq '.alerts[0].startsAt = $dateVal' --arg dateVal "$DATE""T22:39:21.245939463Z" |\
  jq '.alerts[0].endsAt = $endsat' --arg endsat "$DATE""T22:40:51.245939463Z" |\
  curl -X POST http://localhost:8801 -H "Content-Type: application/json" --data-binary @-
done

