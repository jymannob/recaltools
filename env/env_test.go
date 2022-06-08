package env

import (
	"os"
	"testing"
)

func TestStringFromEnvironment(t *testing.T) {
	os.Setenv("ENV_SET", "PASS")

	type args struct {
		envName      string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"env var do not set", args{"ENV_NOT_SET", "FAIL"}, "FAIL"},
		{"env var set to PASS", args{"ENV_SET", "FAIL"}, "PASS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringFromEnvironment(tt.args.envName, tt.args.defaultValue); got != tt.want {
				t.Errorf("StringFromEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirFromEnvironment(t *testing.T) {
	dir, _ := os.Getwd()
	os.Setenv("GOOD_ENV_PATH", dir)
	os.Setenv("BAD_ENV_PATH", "./badenv")
	os.Setenv("ENV_PATH_IS_FILE", "./env.go")

	type args struct {
		envName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"env path exist", args{"GOOD_ENV_PATH"}, dir, false},
		{"env path does not exist", args{"BAD_ENV_PATH"}, "", true},
		{"env path is not directory", args{"ENV_PATH_IS_FILE"}, "", true},
		{"env not set", args{"NOT_EXISTING_ENV"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DirFromEnvironment(tt.args.envName)

			if (err != nil) != tt.wantErr {
				t.Errorf("DirFromEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DirFromEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequiredStringFromEnvironment(t *testing.T) {

	os.Setenv("GOOD_ENV_STRING", "PASS")

	type args struct {
		envName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"env not set", args{"NOT_EXISTING_ENV"}, "", true},
		{"env exist", args{"GOOD_ENV_STRING"}, "PASS", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequiredStringFromEnvironment(tt.args.envName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequiredStringFromEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RequiredStringFromEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoolFromEnvironment(t *testing.T) {

	os.Setenv("GOOD_ENV_BOOL", "TRUE")
	os.Setenv("BAD_ENV_FORMAT", "NOPE")

	type args struct {
		envName      string
		defaultValue bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"good env", args{"GOOD_ENV_BOOL", false}, true},
		{"bad format", args{"BAD_ENV_FORMAT", false}, false},
		{"return default", args{"NOT_EXIST", true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolFromEnvironment(tt.args.envName, tt.args.defaultValue); got != tt.want {
				t.Errorf("BoolFromEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
