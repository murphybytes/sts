package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := summarize(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %q", err)
		os.Exit(1)
	}

}

func summarize(rdr io.Reader, wtr io.Writer) (err error) {
	scanner := bufio.NewScanner(rdr)
	bucket := newBucket(&logWriter{wtr})
	// purge any partial time buckets after we've scanned the whole file
	defer func() {
		if err == nil {
			err = bucket.dump()
		}
	}()

	for scanner.Scan() {
		if err = bucket.add(scanner.Text()); err != nil {
			return err
		}
	}
	err = scanner.Err()
	return err
}
