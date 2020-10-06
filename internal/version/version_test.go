package version

import "testing"

func TestVersionString(t *testing.T) {
	wantVersion := "test version 0.0.1"
	gotVersion := VersionString("test")
	if gotVersion != wantVersion {
		t.Fatalf("unexpected version got %s\n want %s", gotVersion, wantVersion)
	}
}
