package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRelayAllowIsExplicit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "allow")
	if allowed, forced := readRelayAllow(path); len(allowed) != 0 || len(forced) != 0 {
		t.Fatalf("missing file allowed players: %v %v", allowed, forced)
	}
	if err := os.WriteFile(path, []byte("# test pair\n42\n!77 # forced\nTOUS\nnot-a-pid\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	allowed, forced := readRelayAllow(path)
	if !reflect.DeepEqual(allowed, []uint64{42, 77}) || !reflect.DeepEqual(forced, []uint64{77}) {
		t.Fatalf("unexpected allow list: allowed=%v forced=%v", allowed, forced)
	}
}
