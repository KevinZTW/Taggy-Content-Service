package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

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

type Worker struct {
	workerId string
	mapf     func(string, string) []KeyValue
	reducef  func(string, []string) string
}

func MakeWorker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) *Worker {
	w := &Worker{}
	w.workerId = getId()
	w.mapf = mapf
	w.reducef = reducef
	return w
}

func (w *Worker) Run() {

	task := w.getTask()
	fmt.Printf("get the task %+v", task)
	if task.Phase == mapPhase {
		w.doMapTask(task)
	} else {
		w.doReduceTask(task)
	}
	//todo exec map or reduce func base on task state

}

//send RPC request to get task
func (w *Worker) getTask() Task {

	args, reply := TaskArgs{}, TaskReply{}
	args.WorkerId = w.workerId
	fmt.Println("send req to get task")
	call("Coordinator.GetTask", &args, &reply)
	return reply.Task
}

func (w *Worker) doMapTask(t Task) {
	println("do map")
	contents, err := ioutil.ReadFile(t.FileName)
	if err != nil {
		fmt.Println("error to read file")
		return
	}

	kvs := w.mapf(t.FileName, string(contents))
	reduces := make([][]KeyValue, t.NReduce)
	for _, kv := range kvs {
		idx := ihash(kv.Key) % t.NReduce
		reduces[idx] = append(reduces[idx], kv)
	}

	for idx, l := range reduces {
		fileName := reduceName(t.Seq, idx)
		f, err := os.Create(fileName)
		if err != nil {
			w.reportTask(t, false, err)
			return
		}
		enc := json.NewEncoder(f)
		for _, kv := range l {
			if err := enc.Encode(&kv); err != nil {
				w.reportTask(t, false, err)
			}

		}
		if err := f.Close(); err != nil {
			w.reportTask(t, false, err)
		}
	}
	w.reportTask(t, true, nil)

}

func (w *Worker) doReduceTask(t Task) {
	println("do reduce")
	maps := make(map[string][]string)
	for idx := 0; idx < t.NMaps; idx++ {
		fileName := reduceName(idx, t.Seq)
		println(fileName)
		file, err := os.Open(fileName)
		if err != nil {
			w.reportTask(t, false, err)
			return
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			if _, ok := maps[kv.Key]; !ok {
				maps[kv.Key] = make([]string, 0, 100)
			}
			maps[kv.Key] = append(maps[kv.Key], kv.Value)
		}
	}

	res := make([]string, 0, 100)
	for k, v := range maps {
		res = append(res, fmt.Sprintf("%v %v\n", k, w.reducef(k, v)))
	}

	if err := ioutil.WriteFile(mergeName(t.Seq), []byte(strings.Join(res, "")), 0600); err != nil {
		w.reportTask(t, false, err)
	}

	w.reportTask(t, true, nil)
}

func (w *Worker) reportTask(t Task, done bool, err error) {
	if err != nil {
		log.Printf("%v", err)
	}
	args := ReportTaskArgs{}
	args.Done = done
	args.Seq = t.Seq
	args.Phase = t.Phase
	reply := ReportTaskReply{}
	if ok := call("Coordinator.ReportTask", &args, &reply); !ok {
		fmt.Printf("report task fail:%+v", args)
	}
}

//
// start a thread that listens for RPCs from coordinator.go &  worker.go
//
func (w *Worker) server() {
	rpc.Register(w)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func getId() string {

	args, reply := struct{}{}, IdReply{}
	call("Coordinator.GetWorkerId", &args, &reply)
	fmt.Printf("Worker Id is : %s\n", reply.WorkerId)
	return reply.WorkerId
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")

	// sockname := coordinatorSock()
	// c, err := rpc.DialHTTP("unix", sockname)
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
