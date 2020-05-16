package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type record struct {
	pid      int
	taskTime time.Time
	done     bool
	doneTime time.Time
}

// Master hold filenames
type Master struct {
	mu     sync.Mutex
	task   map[string]record
	files  []string
	cursor int
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	fmt.Println("I'm in example ", args.X)
	reply.Y = args.X + 1
	return nil
}

// GetTask for test
func (m *Master) GetTask(args *TaskRequestArgs, reply *TaskRequestReplyArgs) error {
	fmt.Println("I'm in echo ", args.Numbs)
	m.mu.Lock()
	defer m.mu.Unlock()
	end := m.cursor + args.Numbs
	if end > len(m.files) {
		end = len(m.files)
	}
	reply.FileNames = m.files[m.cursor:end]
	m.cursor = end
	//todo add task record
	for _, file := range reply.FileNames {
		m.task[file] = record{args.Pid, time.Now(), false, time.Time{}}
	}
	fmt.Println(m.files)
	fmt.Println(m)
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.
	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}
	m.files = files
	m.task = map[string]record{}
	// Your code here.

	m.server()
	return &m
}
