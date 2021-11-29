package main_test

import (
	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
	Vigilant "vigilant"
)

func TestGetPostJson(t *testing.T) {
	a := assert.New(t)
	b, err := ioutil.ReadFile("example_web_hook_alert_message.json")
	a.NoError(err)

	startsAt, _ := time.Parse(time.RFC3339, "2021-11-22T22:39:21.245939463Z")
	endsAt, _ := time.Parse(time.RFC3339, "2021-11-22T22:40:51.245939463Z")
	expectedStruct := template.Data{
		Receiver:          "web\\.hook",
		Status:            "resolved",
		Alerts:            []template.Alert{
			template.Alert{
				Status:       "resolved",
				Labels:       template.KV{
					// TODO: how to fill in Key-Value?
					"alertname": "uptime_low",
					"instance": "localhost:8099",
					"job": "local_app",
					"severity": "warning",
					"type": "up_time",
				},
				Annotations:  template.KV{
					"summary": "Application uptime low",
				},
				StartsAt:     startsAt,
				EndsAt:       endsAt,
				GeneratorURL: "http://rpmbp.local:9090/graph?g0.expr=beam_up_time+%3E+1\\u0026g0.tab=1",
				Fingerprint:  "8009035229a59d05",
			},
		},
		GroupLabels:       template.KV{
			"alertname": "uptime_low",
		},
		CommonLabels:      template.KV{
			"alertname": "uptime_low",
			"instance": "localhost:8099",
			"job": "local_app",
			"severity": "warning",
			"type": "up_time",
		},
		CommonAnnotations: template.KV{
			"summary": "Application uptime low",
		},
		ExternalURL:       "http://rpmbp.local:9093",
	}

	a.Equal(expectedStruct, Vigilant.GetPostJson(b))
}
