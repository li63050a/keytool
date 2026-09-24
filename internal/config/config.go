package config

import (
	"os"
	"path/filepath"
	"strings"
)

const AppName = "keytool"

func HomeDir() string {
	if env := os.Getenv("KEYTOOL_HOME"); env != "" {
		return env
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".li", AppName)
}

func KeysDir() string   { return filepath.Join(HomeDir(), "keys") }
func GPGDir() string    { return filepath.Join(KeysDir(), "gpg") }
func SSHDir() string    { return filepath.Join(KeysDir(), "ssh") }
func CertDir() string   { return filepath.Join(HomeDir(), "certs") }
func CADir() string     { return filepath.Join(HomeDir(), "ca") }
func CsrDir() string    { return filepath.Join(HomeDir(), "csr") }
func BackupDir() string { return filepath.Join(HomeDir(), "backup") }
func LangFile() string  { return filepath.Join(HomeDir(), ".lang") }

func EnsureDirs() error {
	for _, d := range []string{
		HomeDir(), KeysDir(), GPGDir(), SSHDir(),
		CertDir(), CADir(), CsrDir(), BackupDir(),
	} {
		if err := os.MkdirAll(d, 0700); err != nil {
			return err
		}
	}
	return nil
}

func SaveLang(l string) error { return os.WriteFile(LangFile(), []byte(l), 0600) }

func LoadLang() string {
	data, err := os.ReadFile(LangFile())
	if err != nil {
		return "zh"
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		return "zh"
	}
	return s
}
