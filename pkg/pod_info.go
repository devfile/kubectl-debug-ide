package pkg

type ContainerEndpoint struct {
	name       string
	targetPort int
}

type ContainerEnv struct {
	name  string
	value string
}

type ContainerVolume struct {
	name string
	path string
}

type ContainerInfo struct {
	name          string
	image         string
	command       []string
	args          []string
	env           []ContainerEnv
	endpoints     []ContainerEndpoint
	volumes       []ContainerVolume
	memoryRequest string
	memoryLimit   string
	cpuRequest    string
	cpuLimit      string
}

type PodInfo struct {
	containers []ContainerInfo
}
