package main

import (
	"log"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/zencoder/ddbsync"
)

var (
	table   = kingpin.Arg("table", "Name of DynamoDB table used for locking.").Required().String()
	mutex   = kingpin.Arg("name", "Name of the lock to acquire.").Required().String()
	ttl     = kingpin.Flag("ttl", "TTL to create lock with.").Default("10m").Duration()
	timeout = kingpin.Flag("timeout", "Time out and fail if the lock isn't acquired within this duration.").Default("10m").Duration()
	region  = kingpin.Flag("region", "AWS region in which to use DynamoDB").Default("us-east-1").String()
	unlock  = kingpin.Flag("unlock", "Unlock the named mutex.").Short('u').Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	m := ddbsync.NewMutex(
		*mutex,
		int64(*ttl),
		ddbsync.NewDatabase(*table, *region, "", false),
	)

	if *unlock {
		m.Unlock()
		return
	}

	acquired := make(chan struct{})
	timer := time.After(*timeout)

	go func() {
		m.Lock()
		acquired <- struct{}{}
	}()

	select {
	case <-acquired:
		return
	case <-timer:
		log.Fatalf("failed to acquire lock within timeout window")
	}
}
