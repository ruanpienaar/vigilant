#!/bin/bash

for i in {0..10}
do
  cat example_web_hook_alert_message.json |\
  jq '.alerts[0].labels.job = $newVal' --arg newVal "job-$i" |\
  curl -X POST http://localhost:8801 -H "Content-Type: application/json" --data-binary @-
done

