package main

import "testing"

func TestStoreReturnsCopy(t *testing.T) {
	store := NewStore()

	store.Set([]CheckResult{
		{URL: "http://example.local", Online: true},
	})

	first := store.Latest()
	first[0].URL = "changed"

	second := store.Latest()

	if second[0].URL != "http://example.local" {
		t.Errorf("store data was modified externally: %+v", second[0])
	}
}

func TestStoreSnapshot(t *testing.T) {
	store := NewStore()

	if snap := store.Snapshot(); snap.UpdatedAt != nil || snap.Summary.Total != 0 {
		t.Errorf("empty store snapshot = %+v", snap)
	}

	store.Set([]CheckResult{
		{URL: "a", Online: true},
		{URL: "b", Online: false},
		{URL: "c", Online: true},
	})

	snap := store.Snapshot()

	if snap.UpdatedAt == nil {
		t.Error("UpdatedAt should be set after Set")
	}

	want := Summary{Total: 3, Online: 2, Offline: 1}
	if snap.Summary != want {
		t.Errorf("Summary = %+v; want %+v", snap.Summary, want)
	}
}
