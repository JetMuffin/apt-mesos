package core

import (
	"fmt"
	
	"github.com/golang/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
)
type Volume struct {
	ContainerPath string `json:"container_path,omitempty"`
	HostPath      string `json:"host_path,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

type Task struct {
	ID      string
	Command []string
	Image   string
	Volumes []*Volume
}

func createTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, task *Task) *mesosproto.TaskInfo {
	taskInfo := mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("volt-task-%s", task.ID)),
		TaskId: &mesosproto.TaskID{
			Value: &task.ID,
		},
		SlaveId:   offer.SlaveId,
		Resources: resources,
		Command:   &mesosproto.CommandInfo{},
	}

	// Set value only if provided
	if task.Command[0] != "" {
		taskInfo.Command.Value = &task.Command[0]
	}

	// Set args only if they exist
	if len(task.Command) > 1 {
		taskInfo.Command.Arguments = task.Command[1:]
	}

	// Set the docker image if specified
	if task.Image != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &task.Image,
			},
		}

		for _, v := range task.Volumes {
			var (
				vv   = v
				mode = mesosproto.Volume_RW
			)

			if vv.Mode == "ro" {
				mode = mesosproto.Volume_RO
			}

			taskInfo.Container.Volumes = append(taskInfo.Container.Volumes, &mesosproto.Volume{
				ContainerPath: &vv.ContainerPath,
				HostPath:      &vv.HostPath,
				Mode:          &mode,
			})
		}

		taskInfo.Command.Shell = proto.Bool(false)
	}

	return &taskInfo
}