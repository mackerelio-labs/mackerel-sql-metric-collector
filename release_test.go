package collector

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"gopkg.in/yaml.v2"
)

// This file contains test cases to check name_template in .goreleaser.yaml generates expected filename.

type releaseConfig struct {
	Archives []struct {
		NameTemplate string `yaml:"name_template"`
	} `yaml:"archives"`
}

type buildArg struct {
	ProjectName string
	Os          string
	Arch        string
	Arm         int
}

func TestReleaseName(t *testing.T) {
	text, err := os.ReadFile(".goreleaser.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var cfg releaseConfig
	if err := yaml.Unmarshal(text, &cfg); err != nil {
		t.Fatal(err)
	}
	tmpl, err := template.New("name").Parse(cfg.Archives[0].NameTemplate)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]buildArg{
		"pkg_linux_x86_64":   {"pkg", "linux", "amd64", 0},
		"pkg_linux_i386":     {"pkg", "linux", "386", 0},
		"pkg_windows_x86_64": {"pkg", "windows", "amd64", 0},
		"pkg_darwin_arm64":   {"pkg", "darwin", "arm64", 0},
		"pkg_darwin_armv6":   {"pkg", "darwin", "arm", 6},
	}
	for name, arg := range tests {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, arg); err != nil {
			t.Fatal(err)
		}
		if s := buf.String(); s != name {
			t.Errorf("Execute(%v) = %s; want %s", arg, s, name)
		}
	}
}
