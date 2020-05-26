package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BatchSize task file size
const BatchSize = 3

// PathIntermediate mr intermediate output directory
const PathIntermediate = "intermediate"

// ByKey for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	if _, err := os.Stat(PathIntermediate); os.IsNotExist(err) {
		//could use once mutex
		os.Mkdir(PathIntermediate, 0700)
	}
	// Your worker implementation here.

	// uncomment to send the Example RPC to the master.
	for {
		task := requestTask(BatchSize)

		if task.TaskType == TaskMapType {
			fmt.Println("task id ", task.TaskID)
			if len(task.FileNames) == 0 {
				fmt.Println("no map task get, I am ready to sleep for a while!")
				time.Sleep(time.Second * 5)
				continue
			}
			err := processMap(mapf, task)
			if err != nil {
				log.Fatalf("map failed %s", err)
			} else {
				updateMapTaskSuccess(task.TaskID)
			}

		} else if task.TaskType == TaskReduceType {
			log.Println("let me see reduce bucket is", task.ReduceBucket)
			log.Println("I have no  Idea for task", task.TaskID)
			err := processReduce(reducef, task)
			if err != nil {
				log.Fatalf("error in process Reduce %s\n", err)
			}
			time.Sleep(time.Second * 5)
		} else {
			log.Println("May be we should wait!")
			time.Sleep(time.Second * 5)
		}
	}

	// ready for reduce

}

func processMap(mapf func(string, string) []KeyValue, task TaskRequestReplyArgs) error {
	buckets := make([][]KeyValue, task.ReduceNum)
	for _, filename := range task.FileNames {
		fmt.Println("request filename:", filename)
		fmt.Println("Task id  is:", task.TaskID)
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open map file %v", filename)
			return err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
			return err
		}
		file.Close()
		kva := mapf(filename, string(content))
		// intermediate = append(intermediate, kva...)
		for _, kv := range kva {
			buckets[ihash(kv.Key)%task.ReduceNum] = append(buckets[ihash(kv.Key)%task.ReduceNum], kv)
		}
	}

	for i, bucket := range buckets {
		iterName := fmt.Sprintf("%s/mr-%d-%d", PathIntermediate, task.TaskID, i)
		ofile, err := os.Create(iterName)
		// fmt.Fprintf(ofile, "abc")
		if err != nil {
			//todo task error
			log.Fatal(fmt.Sprintf("create file error %s", err))
			return err
		}
		enc := json.NewEncoder(ofile)
		for _, kv := range bucket {
			err := enc.Encode(&kv)
			if err != nil {
				//todo task error
				log.Fatal(fmt.Sprintf("create file error %s", err))
				return err
			}
		}
		ofile.Close()
	}
	return nil
}

func processReduce(reducef func(string, []string) string, task TaskRequestReplyArgs) error {
	files, err := getReduceFiles(task.ReduceBucket)
	if err != nil {
		log.Fatalf("cannot cat reduce files %v", err)
	}
	fmt.Printf("reducefile %#v\n", files)
	intermediate := []KeyValue{}

	for _, f := range files {
		file, err := os.Open(PathIntermediate + "/" + f)
		if err != nil {
			log.Printf("cannot open inter %v", f)
			return err
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
	}

	sort.Sort(ByKey(intermediate))

	oname := fmt.Sprintf("%s/mr-out-%d", PathIntermediate, task.ReduceBucket)
	ofile, _ := os.Create(oname)

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()

	return nil
}

func getReduceFiles(bucket int) ([]string, error) {
	currentfiles, err := ioutil.ReadDir(PathIntermediate)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	reduceFiles := []string{}
	for _, file := range currentfiles {
		f := file.Name()
		if strings.HasPrefix(f, "mr") {
			split := strings.Split(f, "-")
			if split[2] == strconv.Itoa(bucket) {
				reduceFiles = append(reduceFiles, f)
			}
		}
	}
	return reduceFiles, nil
}

func requestTask(nums int) TaskRequestReplyArgs {
	pid := os.Getpid()
	args := TaskRequestArgs{nums, pid}
	reply := TaskRequestReplyArgs{}
	call("Master.GetTask", &args, &reply)
	fmt.Println(reply.Err)
	if reply.Err != "" {
		fmt.Println("have error ", reply.Err)
	}
	return reply
}

func updateMapTaskSuccess(taskId int) UpdateStatusReply {
	args := UpdateStatusRequest{"M", taskId, SUCCESS}
	reply := UpdateStatusReply{}
	call("Master.UpdateMapTaskStatus", &args, &reply)
	fmt.Println(reply.Err)
	if reply.Err != "" {
		fmt.Println("have error ", reply.Err)
	}
	return reply
}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// call("Master.Example", &args, &reply)
	call("Master.Echo", &args, &reply)
	// reply.Y should be 100.
	fmt.Printf("reply.Y %v\n", reply.Y)
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
