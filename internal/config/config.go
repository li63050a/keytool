package config

import (
	"os"
	"path/filepath"
)

const AppName = "keytool"

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "."+AppName)
}

func KeysDir() string   { return filepath.Join(HomeDir(), "keys") }
func GPGDir() string    { return filepath.Join(KeysDir(), "gpg") }
func SSHDir() string    { return filepath.Join(KeysDir(), "ssh") }

func EnsureDirs() error {
	for _, d := range []string{HomeDir(), KeysDir(), GPGDir(), SSHDir()} {
		if err := os.MkdirAll(d, 0700); err != nil {
			return err
		}
	}
	return nil
}
