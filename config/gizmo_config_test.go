package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Logf("TestLoadConfig")

	filename := "./examples/gizmo_config.json"
	LoadConfig(filename)
}
