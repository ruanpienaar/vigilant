package main_test

import (
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ruanpienaar/vigilant"
)

func TestGetPostJson(t *testing.T) {
	a := assert.New(t)
	b, err := ioutil.ReadFile("example_web_hook_alert_message.json")
	a.NoError(err)
	Vigilant.GetPostJson(b)
}
