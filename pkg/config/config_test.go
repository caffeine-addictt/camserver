package config_test

import (
	"net/url"
	"testing"

	"github.com/caffeine-addictt/camserver/pkg/config"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestRTSP(t *testing.T) {
	tt := []struct {
		str    string
		err    bool
		expect config.Rtsp
	}{
		{str: "rtsp://example.com", err: false, expect: config.Rtsp{URL: &url.URL{Scheme: "rtsp", Host: "example.com:554"}}},
		{str: "rtsp://example.com:2", err: false, expect: config.Rtsp{URL: &url.URL{Scheme: "rtsp", Host: "example.com:2"}}},
		{str: "rtsps://example.com", err: false, expect: config.Rtsp{URL: &url.URL{Scheme: "rtsps", Host: "example.com:554"}}},
		{str: "https://example.com:2", err: true, expect: config.Rtsp{}},
	}

	for _, tc := range tt {
		var out config.Rtsp
		err := yaml.Unmarshal([]byte(tc.str), &out)

		if tc.err {
			assert.Error(t, err, "expected error")
			continue
		}
		if assert.NoError(t, err, "expected no error") {
			assert.Equal(t, tc.expect.URL, out.URL)
		}
	}
}
