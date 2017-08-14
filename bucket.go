package main

import (
	"time"
)

type bucket struct {
	currentBucket time.Time
	out           outputter
	entries       []logEntry
}

const bucketDuration = time.Minute

func newBucket(wtr outputter) *bucket {
	return &bucket{out: wtr}
}

func (b *bucket) add(line string) error {
	entry, err := newLogEntry(line)
	if err != nil {
		return err
	}
	if b.currentBucket.IsZero() {
		b.currentBucket = entry.timeStamp.Truncate(bucketDuration)
		b.entries = append(b.entries, *entry)
		return nil
	}
	entryBucket := entry.timeStamp.Truncate(bucketDuration)
	if entryBucket.After(b.currentBucket) {
		b.out.output(b.currentBucket, b.entries)
		b.reset()
		b.currentBucket = entryBucket
		b.entries = append(b.entries, *entry)
	} else {
		b.entries = append(b.entries, *entry)
	}
	return nil
}

func (b *bucket) dump() error {
	b.out.output(b.currentBucket, b.entries)
	b.reset()
	return nil
}

func (b *bucket) reset() {
	lookup := map[string]struct{}{}
	var newEnts []logEntry

	for _, entry := range b.entries {
		if _, ok := lookup[entry.host]; !ok {
			lookup[entry.host] = struct{}{}
			newEnts = append(newEnts, logEntry{host: entry.host})
		}
	}

	b.entries = newEnts

}
