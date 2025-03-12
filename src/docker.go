package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"archive/tar"
	"bytes"
	"io"
	"sort"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var docker *client.Client

type ImageMeta struct {
	Name     string
	Type     string
	Template string
}

type ContainerSetup struct {
	Name       string
	Image      string
	StartupCmd string
	Volumes    []string
}

type ContainerState string

const (
	ContainerStateStart ContainerState = "start"
	ContainerStateStop  ContainerState = "stop"
)

func dockerConnect() {
	var err error
	docker, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	version := docker.ClientVersion()

	if version > "1.47" {
		log.Fatalf("Docker client version is too new: %s", version)
	}
}

func ListImages(meta ImageMeta) ([]image.Summary, error) {
	filterArgs := buildFilterArgs(meta)

	images, err := docker.ImageList(context.Background(), image.ListOptions{
		All:     false,
		Filters: filterArgs,
	})

	if err != nil {
		log.Fatalf("Error listing Docker images in ListImages: %v", err)
		return nil, err
	}

	return images, nil
}

func ListContainers() ([]container.Summary, error) {
	containers, err := docker.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", "source="+appName), filters.Arg("label", "type=user")),
	})

	if err != nil {
		log.Fatalf("Error listing Docker containers: %v", err)
		return nil, err
	}

	return containers, nil
}

func StartContainer(id string) error {
	err := docker.ContainerStart(context.Background(), id, container.StartOptions{})
	return err
}

func StopContainer(id string) error {
	err := docker.ContainerStop(context.Background(), id, container.StopOptions{})
	return err
}

func DeleteContainer(id string) error {
	err := docker.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})

	return err
}

func FindContainersWithImage(tag string) ([]container.Summary, error) {
	containers, err := docker.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("ancestor", tag)),
	})

	if err != nil {
		log.Println("Error listing Docker containers:", err)
		return nil, err
	}

	return containers, nil
}

func FindImagesByMeta(meta ImageMeta) ([]image.Summary, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "source="+appName)

	if meta.Name != "" {
		filterArgs.Add("label", "name="+meta.Name)
	}

	if meta.Template != "" {
		filterArgs.Add("label", "template="+meta.Template)
	}

	images, err := docker.ImageList(context.Background(), image.ListOptions{
		All:     true,
		Filters: filterArgs,
	})

	if err != nil {
		log.Fatalf("Error listing Docker images in FindImagesByMeta: %v", err)
		return nil, err
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Created > images[j].Created
	})

	return images, nil
}

func FindImageById(id string) (*image.InspectResponse, error) {
	image, err := docker.ImageInspect(context.Background(), id)

	if err != nil {
		log.Println("Error inspecting Docker image", err)
		return nil, err
	}

	return &image, err
}

func ImageIdExists(id string) bool {
	images, err := docker.ImageList(context.Background(), image.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("id", id)),
	})

	if err != nil {
		return false
	}

	return len(images) > 0
}

func DeleteImage(id string) error {
	imageToDelete := id
	name := ""
	containers, err := FindContainersWithImage(id)

	if err != nil {
		log.Println("Error finding containers with image", err)
		return err
	}

	log.Println("found:", len(containers), "containers with image:", id)

	for _, container := range containers {
		log.Println("Deleting container:", container.ID)
		err := DeleteContainer(container.ID)

		if err != nil {
			log.Println("Error deleting container:", container.ID, err)
			return err
		}
	}

	for true {
		_image, err := FindImageById(imageToDelete)

		if err != nil {
			log.Println("Error finding image:", imageToDelete, err)
			return err
		}

		log.Println("Deleting image:", _image.ID, _image.Config.Labels)

		if len(name) == 0 {
			name = _image.Config.Labels["name"]
		}

		_, err = docker.ImageRemove(context.Background(), _image.ID, image.RemoveOptions{
			Force:         true,
			PruneChildren: true,
		})

		if err != nil {
			log.Println("Error removing Docker image:", err)
			return err
		}

		if _image.Config.Labels["name"] != name {
			log.Println("Skipping system image:", _image.Config.Labels["name"], _image.Parent)
			break
		}

		imageToDelete = _image.Parent
	}

	return nil
}

func SetContainerState(id string, state ContainerState) error {
	if state == ContainerStateStart {
		return StartContainer(id)
	}

	if state == ContainerStateStop {
		return StopContainer(id)
	}

	return nil
}

func FindAvailablePort(min int, max int) (int, error) {
	containers, err := ListContainers()

	if err != nil {
		log.Fatalf("Error listing Docker containers: %v", err)
		return 0, err
	}

	ports := make(map[int]bool)

	for _, container := range containers {
		i, err := strconv.Atoi(container.Labels["port"])

		if err != nil {
			log.Println("Error parsing port:", err)
			continue
		}

		ports[i] = true
	}

	for i := min; i <= max; i++ {
		if _, ok := ports[i]; !ok {
			return i, nil
		}
	}

	return 0, fmt.Errorf("No available port found in range %d-%d", min, max)
}

func CreateContainer(setup ContainerSetup) (string, error) {
	_ = docker.ContainerRemove(context.Background(), setup.Name, container.RemoveOptions{
		Force: true,
	})

	password := randomString(24)
	availablePort, err := FindAvailablePort(14500, 14599)

	if err != nil {
		log.Fatalf("Error finding available port: %v", err)
		return "", err
	}

	port := strconv.Itoa(availablePort)
	startupCmd := fmt.Sprintf(startupCmdTemplate, setup.StartupCmd, port, password)
	log.Println("Creating container with password:", password, "and port:", availablePort, "using command:", startupCmd, "name", setup.Name)

	exposedPorts := nat.PortSet{}
	exposedPorts[nat.Port(port+"/tcp")] = struct{}{}
	portBindings := nat.PortMap{}

	portBindings[nat.Port(port+"/tcp")] = []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: port,
		},
		{
			HostIP:   "::",
			HostPort: port,
		},
	}

	response, err := docker.ContainerCreate(context.Background(), &container.Config{
		Image: setup.Image,
		Cmd:   parseBashString(startupCmd),
		Labels: map[string]string{
			"source":   appName,
			"type":     "user",
			"name":     setup.Name,
			"template": setup.Image,
			"password": password,
			"port":     port,
		},
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		Binds:        setup.Volumes,
		PortBindings: portBindings,
	}, &network.NetworkingConfig{}, &v1.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}, strings.ToLower(appName)+"_"+setup.Name)

	if err != nil {
		log.Println("Error creating Docker container", err)
		return "", err
	}

	return response.ID, nil
}

func buildFilterArgs(meta ImageMeta) filters.Args {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "source="+appName)

	if meta.Name != "" {
		filterArgs.Add("label", "name="+meta.Name)
	}

	if meta.Type != "" {
		filterArgs.Add("label", "type="+meta.Type)
	}

	if meta.Template != "" {
		filterArgs.Add("label", "template="+meta.Template)
	}

	return filterArgs
}

func ImageExists(meta ImageMeta) bool {
	filterArgs := buildFilterArgs(meta)

	images, err := docker.ImageList(context.Background(), image.ListOptions{
		Filters: filterArgs,
	})

	if err != nil {
		log.Println("Error listing Docker images in ImageExists:", err)
	}

	return len(images) > 0
}

func CreateImage(dockerfile string, meta ImageMeta) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	dockerfileBytes := []byte(dockerfile)

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(dockerfileBytes)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatalf("Error writing tar header: %v", err)
		return nil, err
	}

	if _, err := tw.Write(dockerfileBytes); err != nil {
		log.Fatalf("Error writing Dockerfile content: %v", err)
		return nil, err
	}

	if err := tw.Close(); err != nil {
		log.Fatalf("Error closing tar writer: %v", err)
		return nil, err
	}

	ctx := context.Background()

	imageBuildResponse, err := docker.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			Tags:   []string{strings.ToLower(appName) + "_" + meta.Name + ":latest"},
			Labels: map[string]string{"source": appName, "type": meta.Type, "name": meta.Name, "template": meta.Template},
		},
	)

	if err != nil {
		log.Fatalf("Error building Docker image: %v", err)
		return nil, err
	}

	return imageBuildResponse.Body, nil
}
