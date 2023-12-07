package genenv

type DockerCompose struct {
	Services map[string]map[string]interface{} `yaml:"services"`
}

type Service struct {
	Image         string            `yaml:"image"`
	Build         Build             `yaml:"build"`
	Volumes       []string          `yaml:"volumes"`
	Ports         []string          `yaml:"ports"`
	Environment   map[string]string `yaml:"environment"`
	Command       []string          `yaml:"command"`
	ContainerName string            `yaml:"container_name"`
	Networks      []string          `yaml:"networks"`
}

type Build struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}

func NewDockerCompose() *DockerCompose {
	servicesMap := make(map[string]map[string]interface{})
	return &DockerCompose{
		Services: servicesMap,
	}
}
