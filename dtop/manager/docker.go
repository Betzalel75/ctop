package manager

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/Betzalel75/ctop/dtop/resource"
	api "github.com/fsouza/go-dockerclient"
)

type DockerResourceManager struct {
	client *api.Client
}

func NewDockerResourceManager(client *api.Client) *DockerResourceManager {
	return &DockerResourceManager{client: client}
}

func (drm *DockerResourceManager) LoadContainers() ([]list.Item, error) {
	containers, err := drm.client.ListContainers(api.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, c := range containers {
		name := strings.TrimPrefix(c.Names[0], "/")
		desc := fmt.Sprintf("ID: %s | Status: %s", c.ID[:12], c.Status)
		items = append(items, resource.ResourceItem{
			Id:    c.ID,
			Title: name,
			Desc:  desc,
		})
	}
	return items, nil
}

func (drm *DockerResourceManager) LoadImages() ([]list.Item, error) {
	images, err := drm.client.ListImages(api.ListImagesOptions{All: true})
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, img := range images {
		tag := "<none>:<none>"
		if len(img.RepoTags) > 0 {
			tag = img.RepoTags[0]
		}
		desc := fmt.Sprintf("ID: %s | Size: %dMB", img.ID[7:19], img.Size/1e6)
		items = append(items, resource.ResourceItem{
			Id:    img.ID,
			Title: tag,
			Desc:  desc,
		})
	}
	return items, nil
}

func (drm *DockerResourceManager) LoadVolumes() ([]list.Item, error) {
	volumes, err := drm.client.ListVolumes(api.ListVolumesOptions{})
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, v := range volumes {
		items = append(items, resource.ResourceItem{
			Id:    v.Name,
			Title: v.Name,
			Desc:  "Driver: " + v.Driver,
		})
	}
	return items, nil
}

func (drm *DockerResourceManager) DeleteContainers(items []resource.ResourceItem) []error {
	var errs []error
	for _, item := range items {
		err := drm.client.RemoveContainer(api.RemoveContainerOptions{
			ID:    item.Id,
			Force: true,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete container %s: %w", item.Title, err))
		}
	}
	return errs
}

func (drm *DockerResourceManager) DeleteImages(items []resource.ResourceItem) []error {
	var errs []error
	for _, item := range items {
		err := drm.client.RemoveImageExtended(item.Id, api.RemoveImageOptions{Force: true})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete image %s: %w", item.Title, err))
		}
	}
	return errs
}

func (drm *DockerResourceManager) DeleteVolumes(items []resource.ResourceItem) []error {
	var errs []error
	for _, item := range items {
		err := drm.client.RemoveVolumeWithOptions(api.RemoveVolumeOptions{Name: item.Id})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete volume %s: %w", item.Title, err))
		}
	}
	return errs
}

func (drm *DockerResourceManager) PublishImages(items []resource.ResourceItem, registry, username, tag string) []error {
	var errs []error
	for _, item := range items {
		baseImageName := strings.Split(item.Title, ":")[0]
		targetImage := fmt.Sprintf("%s/%s/%s:%s", registry, username, baseImageName, tag)

		err := drm.client.TagImage(item.Id, api.TagImageOptions{
			Repo:  targetImage,
			Force: true,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to tag %s: %w", item.Title, err))
			continue
		}

		err = drm.client.PushImage(api.PushImageOptions{
			Name:         targetImage,
			Tag:          tag,
			OutputStream: os.Stdout,
			Context:      context.Background(),
		}, api.AuthConfiguration{})
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to push %s: %w", targetImage, err))
		}
	}
	return errs
}