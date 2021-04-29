package agent

var docker *DockerAgent

func init() {
	var err error
	docker, err = NewDockerAgent()
	if err != nil {
		panic(err)
	}
}

func main() {
	docker.Run()
}
