package web

import (
	config "github.com/whoismissing/gizmo/config"

	"testing"
)

func TestGenerateScoreboardHTML(t *testing.T) {
	t.Logf("TestGenerateScoreboardHTML")

	filename := "../config/examples/gizmo_config.json"
	teams := config.LoadConfig(filename)

	t.Logf("+%v", teams)

	_ = GenerateScoreboardHTML(teams)
}
