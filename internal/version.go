package internal

import "encoding/json"

var version string
var revision string

// GetVersions returns the serialized version and revision information for this exporter.
func GetVersions() string {
	m := map[string]string{
		"version":  version,
		"revision": revision,
	}
	v, _ := json.Marshal(m)
	return string(v)
}
