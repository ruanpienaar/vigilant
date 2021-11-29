#!/bin/bash

cat example_web_hook_alert_message.json | curl -X POST http://192.168.0.37:8801 -H "Content-Type: application/json" --data-binary @-