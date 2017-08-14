package main

import (
	"fmt"
	"io"
	"sort"
	"time"
)

type outputter interface {
	output(ts time.Time, entries []logEntry)
}

type logWriter struct {
	wtr io.Writer
}

type summaries []*summary

func (s summaries) Len() int { return len(s) }
func (s summaries) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s summaries) Less(i, j int) bool {
	return s[i].host < s[j].host
}

func (lw *logWriter) output(ts time.Time, entries []logEntry) {
	seen := map[string]int{}
	empty := map[string]struct{}{}
	var sums summaries
	for _, entry := range entries {
		if entry.serviceTime == 0 {
			empty[entry.host] = struct{}{}
			continue
		}
		idx, ok := seen[entry.host]
		if !ok {
			sums = append(sums, newSummary(ts, entry))
			seen[entry.host] = len(sums) - 1
		} else {
			sums[idx].record(entry)
		}
	}
	for k, _ := range empty {
		if _, ok := seen[k]; !ok {
			sums = append(sums, newSummary(ts, logEntry{host: k}))
		}
	}
	sort.Sort(sums)
	for _, s := range sums {
		fmt.Fprint(lw.wtr, s)
	}
}

type summary struct {
	ts    time.Time
	host  string
	count int
	total int
	min   int
	max   int
}

func newSummary(ts time.Time, ent logEntry) *summary {
	count := 0
	if ent.serviceTime > 0 {
		count = 1
	}
	return &summary{
		ts:    ts,
		host:  ent.host,
		count: count,
		min:   ent.serviceTime,
		max:   ent.serviceTime,
		total: ent.serviceTime,
	}
}

func (s *summary) record(ent logEntry) {
	s.count++
	s.total += ent.serviceTime
	if s.min > ent.serviceTime {
		s.min = ent.serviceTime
	}
	if s.max < ent.serviceTime {
		s.max = ent.serviceTime
	}
}

func (s *summary) String() string {
	return fmt.Sprintf("%s,%s,%d,%d,%d,%d\n", s.ts.Format("2006-01-02T15:04:05"), s.host, s.count, s.total, s.min, s.max)
}
