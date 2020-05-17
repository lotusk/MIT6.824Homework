package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//
type ExampleArgs struct {
	X int
}

// ExampleReply rpc reply
type ExampleReply struct {
	Y int
}

// TaskRequestArgs job request args
type TaskRequestArgs struct {
	Numbs int
	Pid   int
}

// TaskRequestReplyArgs job request response
type TaskRequestReplyArgs struct {
	FileNames []string
	TaskId    int
	Err       string
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}