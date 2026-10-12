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
