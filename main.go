package main

import (
	"log"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/SmartThingsOSS/ddbsync"
)

var (
	table   = kingpin.Arg("table", "Name of DynamoDB table used for locking.").Required().String()
	mutex   = kingpin.Arg("name", "Name of the lock to acquire/unlock.").Required().String()
	ttl     = kingpin.Flag("ttl", "TTL of the lock.").Default("10m").Duration()
	timeout = kingpin.Flag("timeout", "Time out and fail if the lock isn't acquired within this duration.").Default("10m").Duration()
	region  = kingpin.Flag("region", "AWS region in which to use DynamoDB").Default("us-east-1").String()
	unlock  = kingpin.Flag("unlock", "Unlock the named mutex.").Short('u').Bool()
)

func main() {
	kingpin.Version("0.0.2")
	kingpin.Parse()

	m := ddbsync.NewMutex(
		*mutex,
		int64(ttl.Seconds()),
		ddbsync.NewDatabase(*table, *region, "", false),
	)

	var action func()
	if *unlock {
		action = m.Unlock
	} else {
		action = m.Lock
	}

	done := make(chan struct{})
	timer := time.After(*timeout)

	go func() {
		action()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return
	case <-timer:
		log.Fatalf("failed to complete lock action within timeout window")
	}
}
