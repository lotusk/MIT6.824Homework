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

// TaskMapType type map
const TaskMapType = "M"

// TaskReduceType  type reduce
const TaskReduceType = "R"

// TaskStatus  refer UNASSIGN, ASSIGN,SUCCESS,FAILED,TIMEOUT
type TaskStatus int

const (
	//UNASSIGN task not assign yet
	UNASSIGN TaskStatus = iota
	ASSIGN
	SUCCESS
	FAILED
	TIMEOUT
)

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
	FileNames    []string
	TaskID       int
	ReduceNum    int
	Err          string
	TaskType     string
	ReduceBucket string
}

// UpdateStatusRequest update status when success or failed
type UpdateStatusRequest struct {
	TaskType string
	TaskID   int
	Status   TaskStatus
}

type UpdateStatusReply struct {
	Err string
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
