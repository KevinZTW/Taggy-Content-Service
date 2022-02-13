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

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type TaskArgs struct {
	WorkerId string
}

type TaskReply struct {
	Task Task
}

type IdArgs struct {
	Domain string
}

type IdReply struct {
	WorkerId string
}

type ReportTaskArgs struct {
	Done     bool
	Seq      int
	Phase    mrPhase
	WorkerId int
}

type ReportTaskReply struct {
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/taggy-mr-coordinator-"
	s += strconv.Itoa(os.Getuid())
	return s
}

func workerSock(workerId int) string {
	s := "/var/tmp/taggy-mr-worker-"
	s += strconv.Itoa(os.Getuid())
	return s
}
