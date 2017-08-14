package main

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	var sums summaries
	sums = append(sums, &summary{host: "zzz"}, &summary{host: "bbbb"}, &summary{host: "cccc"})
	sort.Sort(sums)
	if sums[0].host != "bbbb" {
		t.Fatalf("expected bbb got %q", sums[0].host)
	}
	if sums[2].host != "zzz" {
		t.Fatalf("expected zzz got %q", sums[2].host)
	}
}
