package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

type Coordinator struct {
	mu          sync.Mutex
	workerSeq   int
	workerIds   []string
	nReduce     int
	nMap        int
	phase       mrPhase
	mapTasks    []Task
	reduceTasks []Task
}

type Task struct {
	Seq       int
	FileName  string
	WorkerId  string
	Phase     mrPhase
	TaskState TaskState
	NReduce   int
	NMaps     int
}

type mrPhase int

const (
	mapPhase    mrPhase = 0
	reducePhase mrPhase = 1
)

type TaskState int

const (
	taskIdle       TaskState = 0
	taskInProgress TaskState = 1
	taskDone       TaskState = 2
)

//RPC handlers for the worker to call.
func (c *Coordinator) GetWorkerId(args *struct{}, reply *IdReply) error {
	c.workerSeq += 1
	workerId := strconv.Itoa(c.workerSeq)
	c.workerIds = append(c.workerIds, workerId)
	reply.WorkerId = workerId
	return nil
}

func (c *Coordinator) GetTask(args *TaskArgs, reply *TaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("worker : %s req GetTask\n", args.WorkerId)
	fmt.Printf("map task: %d, red task : %d", len(c.mapTasks), len(c.reduceTasks))
	//assign map task if any in idle
	if c.phase == mapPhase {
		for i := 0; i < len(c.mapTasks); i++ {
			mtask := &c.mapTasks[i]
			if mtask.TaskState == taskIdle {
				fmt.Printf("pick %s", mtask.FileName)
				mtask.TaskState = taskInProgress
				reply.Task = *mtask
				return nil
			}
		}
	} else {
		for i := 0; i < len(c.reduceTasks); i++ {
			rtask := &c.reduceTasks[i]
			fmt.Printf("rtask: %+v", rtask)
			if rtask.TaskState == taskIdle {
				fmt.Printf("pick %s", rtask.FileName)
				rtask.TaskState = taskInProgress
				reply.Task = *rtask
				return nil
			}
		}
	}

	return fmt.Errorf("no task idle")
}

func (c *Coordinator) ReportTask(args *ReportTaskArgs, reply *ReportTaskReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("get report task: %+v\n", args)
	if args.Done {
		c.mapTasks[args.Seq].TaskState = taskDone
	}
	if c.mapTaskDone() {
		c.phase = reducePhase
	}

	// if m.taskPhase != args.Phase || args.WorkerId != m.taskStats[args.Seq].WorkerId {
	// 	return nil
	// }

	// if args.Done {
	// 	m.taskStats[args.Seq].Status = TaskStatusFinish
	// } else {
	// 	m.taskStats[args.Seq].Status = TaskStatusErr
	// }

	// go m.schedule()
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	// sockname := coordinatorSock()
	// os.Remove(sockname)
	// l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	mapDone, redDone := c.mapTaskDone(), c.reduceTaskDone()

	if mapDone && redDone {
		return true
	}
	return false
}

func (c *Coordinator) mapTaskDone() bool {

	//check all maptask are done
	for _, task := range c.mapTasks {
		if task.TaskState != taskDone {
			return false
		}
	}
	return true
}

func (c *Coordinator) reduceTaskDone() bool {

	//check all maptask are done
	for _, task := range c.reduceTasks {
		if task.TaskState != taskDone {
			return false
		}
	}
	return true
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.nReduce = nReduce
	c.nMap = 3
	c.workerSeq = 0
	seq := 0

	for _, file := range files {

		mtask := Task{}
		mtask.FileName = file
		mtask.TaskState = taskIdle
		mtask.NReduce = c.nReduce
		mtask.NMaps = c.nMap
		mtask.Seq = seq
		mtask.Phase = mapPhase
		seq++
		c.mapTasks = append(c.mapTasks, mtask)
	}
	seq = 0
	for i := 0; i < c.nReduce; i++ {
		println("add reduce")
		t := Task{}
		t.NMaps = c.nMap
		t.NReduce = c.nReduce
		t.Phase = reducePhase
		t.TaskState = taskIdle
		t.FileName = "reduce"
		t.Seq = seq
		seq++
		c.reduceTasks = append(c.reduceTasks, t)
	}
	c.server()
	return &c
}
