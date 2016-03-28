package registry

// Port mapping of docker container
type Port struct {
	ContainerPort uint32 `json:"container_port,string"`
	HostPort      uint32 `json:"host_port,string"`
}