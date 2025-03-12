package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Init() {
	if !ImageExists(ImageMeta{Name: "baseimage", Type: "system"}) {
		log.Println("Base image not found, creating...")
		stream, err := CreateImage(baseDockerfile, ImageMeta{Name: "baseimage", Type: "system"})

		if err != nil {
			log.Println("Failed to create base image", err)
			return
		}

		handleLogStream(a.ctx, stream)
	}
}

func (a *App) Containers() []container.Summary {
	containers, err := ListContainers()
	runtime.EventsEmit(a.ctx, "syslog")

	if err != nil {
		return []container.Summary{}
	}

	filteredContainers := []container.Summary{}

	for _, container := range containers {
		if len(container.Labels["port"]) > 0 && len(container.Labels["password"]) > 0 {
			filteredContainers = append(filteredContainers, container)
		}
	}

	return filteredContainers
}

func (a *App) CreateTemplate(name string, baseTemplateId string, template string) string {
	log.Println("creating template", name, baseTemplateId, template)
	fullTemplate := "FROM " + baseTemplateId + "\n" + parseTemplate(template)

	stream, err := CreateImage(fullTemplate, ImageMeta{
		Name:     name,
		Type:     "user",
		Template: baseTemplateId,
	})

	if err != nil {
		return err.Error()
	}

	return handleLogStream(a.ctx, stream)
}

func (a *App) CreateContainer(name string, templateId string, startupCmd string, volumes []string) string {
	log.Println("creating container", name, templateId, startupCmd, volumes)

	_, err := CreateContainer(ContainerSetup{
		Name:       name,
		Image:      templateId,
		StartupCmd: startupCmd,
		Volumes:    volumes,
	})

	if err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) DeleteContainer(id string) string {
	log.Println("deleting container", id)
	err := DeleteContainer(id)

	if err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) DeleteTemplate(id string) string {
	log.Println("deleting template", id)
	err := DeleteImage(id)

	if err != nil {
		log.Println("Failed to delete template", err)
		return err.Error()
	}

	log.Println("deleted template", id)
	return ""
}

func (a *App) SetContainerState(id string, state ContainerState) string {
	log.Println("setting container", id, state)
	err := SetContainerState(id, state)

	if err != nil {
		return err.Error()
	}

	return ""
}

func (a *App) GetTemplates(all bool) []image.Summary {
	filters := ImageMeta{
		Type: "user",
	}

	if all {
		filters.Type = ""
	}

	templates, err := ListImages(filters)

	if err != nil {
		log.Println("Failed to get templates", err)
		return []image.Summary{}
	}

	return templates
}

func (a *App) XpraConnect(port string, password string) {
	_ = XpraConnect(port, password)
}

func (a *App) OpenLink(url string) {
	log.Println("opening link", url)
	cmd := exec.Command("xdg-open", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()

	if err != nil {
		log.Fatalf("Failed to open link: %v", err)
	}

	log.Println("opened link", url)
}
