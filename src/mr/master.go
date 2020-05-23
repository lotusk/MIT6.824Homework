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

// type Phase int

// const (
// 	PhaseMap Phase = iota
// 	PhaseReduce
// )

type record struct {
	file     string
	pid      int
	taskTime time.Time
	status   TaskStatus
	doneTime time.Time
	taskID   int
	taskType string
}

// Master hold filenames
type Master struct {
	mu                   sync.Mutex
	record               map[string]*record
	tasks                map[int][]*record
	reduceTasks          map[int]*record
	files                []string
	cursor               int
	taskCursor           int
	nReduce              int
	phase                string
	failedFiles          []string
	successCounter       int
	successReduceCounter int
	reduceBucketCursor   int
}

type argError struct {
	arg  TaskStatus
	prob string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

// Your code here -- RPC handlers for the worker to call.

// Example an example RPC handler.
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
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
			reply.Err = fmt.Sprintf("%s", err)
		}
	}()
	if m.phase == TaskMapType {
		m.getMapTask(args, reply)
	} else if m.phase == TaskReduceType {
		//TODO reduce
		log.Println("let me try to reduce!")
		m.getReduceTask(args, reply)

	} else {
		//TODO done
	}
	return nil
}

// UpdateTaskStatus change task status when task done or error
func (m *Master) UpdateMapTaskStatus(args *UpdateStatusRequest, reply *UpdateStatusReply) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
			reply.Err = fmt.Sprintf("%s", err)
		}
	}()
	m.mu.Lock()
	defer m.mu.Unlock()
	if args.Status == SUCCESS {
		for _, record := range m.tasks[args.TaskID] {
			record.status = SUCCESS
			m.successCounter++
		}

		if m.successCounter >= len(m.files) {
			m.phase = TaskReduceType
			log.Println("All Map task have done , let we entry reduce phase!")
		}
	} else if args.Status == FAILED {
		//TODO add to failed list
	} else {
		return &argError{args.Status, "map task just return success or failed"}
	}
	return nil
}

func (m *Master) getMapTask(args *TaskRequestArgs, reply *TaskRequestReplyArgs) error {
	fmt.Println("I'm in map ", args.Numbs)
	m.mu.Lock()
	defer m.mu.Unlock()
	end := m.cursor + args.Numbs
	if end > len(m.files) {
		end = len(m.files)
	}

	if end > m.cursor {
		replyFiles := m.files[m.cursor:end]
		m.cursor = end
		//todo add task record
		tasks := []*record{}
		for _, file := range replyFiles {
			fmt.Println("put m.task ", file)
			record := &record{file, args.Pid, time.Now(), ASSIGN, time.Time{}, m.taskCursor, TaskMapType}
			m.record[file] = record
			tasks = append(tasks, record)
		}
		m.tasks[m.taskCursor] = tasks
		fmt.Println(m.files)
		for k, v := range m.record {
			fmt.Println(k, v)
		}
		reply.FileNames = replyFiles
		reply.TaskID = m.taskCursor
		reply.ReduceNum = m.nReduce
		reply.TaskType = TaskMapType

		m.taskCursor++
	} else {
		// return empty reply ,let worker  sleep a while
		log.Println("no map task and wait all task done!")
	}
	return nil
}

func (m *Master) getReduceTask(args *TaskRequestArgs, reply *TaskRequestReplyArgs) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
			reply.Err = fmt.Sprintf("%s", err)
		}
	}()

	fmt.Println("I'm in reduce ")
	m.mu.Lock()
	defer m.mu.Unlock()
	//TODO

	reply.TaskType = "R"
	reply.TaskID = m.taskCursor
	reply.ReduceBucket = m.reduceBucketCursor
	m.reduceTasks[m.taskCursor] = &record{"", args.Pid, time.Now(), ASSIGN, time.Time{}, m.taskCursor, TaskReduceType}
	m.reduceBucketCursor++
	m.taskCursor++
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

func start() {

}

// Done check whether it can stop.
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.
	return ret
}

// MakeMaster as the name
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}
	m.files = files
	m.record = map[string]*record{}
	m.tasks = map[int][]*record{}
	m.reduceTasks = map[int]*record{}
	m.nReduce = nReduce
	m.phase = TaskMapType
	// Your code here.

	m.server()
	return &m
}
