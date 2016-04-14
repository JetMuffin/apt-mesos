package registry

import (
	"path"

	"github.com/icsnju/apt-mesos/docker"
	"github.com/icsnju/apt-mesos/utils"
)

type Job struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Image      string             `json:"image"`
	Dockerfile *docker.Dockerfile `json:"dockerfile"`
	ContextDir string             `json:"context_dir"`
	CreateTime int64              `json:"create_time"`
}

func (job *Job) DockerfileExists() bool {
	if job.ContextDir != "" {
		dockerfilePath := path.Join(job.ContextDir, "Dockerfile")
		if !utils.Exists(dockerfilePath) {
			return false
		}
	}
	return true
}
