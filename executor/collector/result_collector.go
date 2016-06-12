package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
	"github.com/mesos/mesos-go/executor"
	"github.com/mesos/mesos-go/mesosproto"
)

// ImageBuilder is the executor to build image
type ImageBuilder struct {
	client *docker.Client
}

// NewImageBuilder return a new image builder
func NewImageBuilder() *ImageBuilder {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	return &ImageBuilder{
		client: client,
	}
}

// Registered called when executor registered to mesos master
func (builder *ImageBuilder) Registered(driver executor.ExecutorDriver, execInfo *mesosproto.ExecutorInfo, fwinfo *mesosproto.FrameworkInfo, slaveInfo *mesosproto.SlaveInfo) {
	fmt.Println("Registered Executor on slave ", slaveInfo.GetHostname())
}

// Reregistered called when executor re-regitered to mesos master
func (builder *ImageBuilder) Reregistered(driver executor.ExecutorDriver, slaveInfo *mesosproto.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ", slaveInfo.GetHostname())
}

// Disconnected called when executor disconnect to mesos master
func (builder *ImageBuilder) Disconnected(executor.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

// LaunchTask called when executor launch tasks
func (builder *ImageBuilder) LaunchTask(driver executor.ExecutorDriver, taskInfo *mesosproto.TaskInfo) {
	fmt.Printf("Launching task %v with ID %v\n", taskInfo.GetName(), taskInfo.GetTaskId().GetValue())

	status := &mesosproto.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesosproto.TaskState_TASK_RUNNING.Enum(),
	}
	_, err := driver.SendStatusUpdate(status)
	if err != nil {
		fmt.Println("Send task running status error: ", err)
	}

	//TODO cp

	fmt.Println("Task finished", taskInfo.GetName())
	status.State = mesosproto.TaskState_TASK_FINISHED.Enum()
	_, err = driver.SendStatusUpdate(status)
	if err != nil {
		fmt.Println("Send task finished status error: ", err)
	}
}

// KillTask called when executor kill tasks
func (builder *ImageBuilder) KillTask(executor.ExecutorDriver, *mesosproto.TaskID) {
	fmt.Println("Kill task")
}

// FrameworkMessage called when executor receive framework's message
func (builder *ImageBuilder) FrameworkMessage(driver executor.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

// Shutdown called when executor shut down
func (builder *ImageBuilder) Shutdown(executor.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
}

// Error called when executor got error
func (builder *ImageBuilder) Error(driver executor.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("Starting Result Collector (APT-MESOS)")
	fmt.Println(os.Args)
	driverConfig := executor.DriverConfig{
		Executor: NewImageBuilder(),
	}
	driver, err := executor.NewMesosExecutorDriver(driverConfig)

	if err != nil {
		fmt.Println("Unable to create a ExecutorDriver ", err.Error())
	}

	_, err = driver.Start()
	if err != nil {
		fmt.Println("Got error:", err)
		return
	}
	fmt.Println("Executor process has started and running.")
	driver.Join()
}