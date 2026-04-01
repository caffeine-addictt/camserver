// Package config
//
// Configuration format for camserver
package config

import (
	"fmt"
	"net"
	"net/url"
)

type Config struct {
	// Archive directory Defaults to `/var/log/camserver/`
	//
	// Relative paths will be resolved relative to the
	// location of the loaded configuration file
	//
	// This the is root directory where all camera recordings
	// will be stored under, in their respective subdirectories
	// corresponding to the SHA256 hash of their URL's `ip`,
	// `port` and `path` (without the credentials)
	//
	// For example:
	//
	//	# camserver.yml
	//	archive_dir: '/var/log/camserver'
	//	cameras:
	//		- name: camera1
	//		  rtsp: 'rtsp://user:password@192.168.1.200:554/Streaming/channels/101'
	//
	// Will use `192.168.1.200:554/Streaming/channels/101` to generate the hash, thus
	// making the final recording directory of `camera1`:
	// `/var/log/camserver/1817dadbf43776ff67f82ba8205bc3bf4e34a5ea03ff65689ab92ea6618bf9e2/`
	ArchiveDirectory string `yaml:"archive_dir,omitempty"`

	// Cameras to look at
	Cameras []CameraCfg `yaml:"cameras"`

	Server ServerCfg `yaml:"server,omitempty"`
}

type CameraCfg struct {
	// Name of the camera
	Name string `yaml:"name"`

	// Camera feed access
	Rtsp Rtsp `yaml:"rtsp"`
}

type Rtsp struct{ *url.URL }

func (r *Rtsp) UnmarshalYAML(data []byte) error {
	u, err := url.Parse(string(data))
	if err != nil {
		return err
	}

	if u.Scheme != "rtsp" && u.Scheme != "rtsps" {
		return fmt.Errorf("invalid scheme %q (expected rtsp/rtsps)", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("rtsp url missing host")
	}

	if u.Port() == "" {
		u.Host = net.JoinHostPort(u.Hostname(), "554")
	}

	*r = Rtsp{URL: u}
	return nil
}

type ServerCfg struct {
	/// Prefer calling [ServerCfg.GetPort]
	Port *uint16 `yaml:"port,omitempty"`
}

func (sc *ServerCfg) GetPort() uint16 {
	if sc.Port == nil {
		return 3000
	}
	return *sc.Port
}
