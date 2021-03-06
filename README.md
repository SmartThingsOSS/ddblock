ddblock (DynamoDB Lock)
=======================

CLI tool to acquire a distributed lock from DynamoDB.

## Building

Install Glide if you don't already have it. On MacOS you can install it with Homebrew: `brew install glide`

To build a binary:
```
glide install
go build
```

To build a linux binary from another platform:
```
glide install
GOOS=linux go build
```

## Setup

Create a DynamoDB table with a primary string key called `Name`.

## Usage

To acquire a lock called "mylock" in a table called "mytable":

`./ddblock mytable mylock`

By default the lock will last for 10 minutes. 

To unlock when you're done:

`./ddblock -u mytable mylock`

To acquire a lock before running a command, run the command, and then release the lock:

```
./ddblock mytable mylock && ( ./my_command ; ./ddblock -u mytable mylock )
```
