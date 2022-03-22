package main

import "os"

var version string = func() string {
	const defaultVer = "0.0.1"
	curVer := os.Getenv("SVC_VERSION")
	if curVer == "" {
		return defaultVer
	}
	return curVer
}()
