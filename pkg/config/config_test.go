package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTruncatesShorterConfig(t *testing.T) {
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	input := "discord:\n  token: test-token\n# trailing comment that must not survive serialization\n"
	for i := 0; i < 100; i++ {
		input += "# padding to make original longer\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	var c Config
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("config.yml")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) >= len(input) {
		t.Fatal("configuration was not truncated")
	}
	var again Config
	if err := again.Load(); err != nil {
		t.Fatalf("second startup: %v", err)
	}
	if again.Discord.Token != "test-token" {
		t.Fatal("token changed")
	}
}
