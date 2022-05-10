package image

import "github.com/docker/docker/api/types"

type ImagePullConfig struct {
	Verbose bool
	All     bool
}

type ImageListConfig struct {
	All bool
}
type ImageRemoveConfig = types.ImageRemoveOptions

type ImageDeletedItem = types.ImageDeleteResponseItem

type ImageRemoveResponse struct {
	DeletedItems []ImageDeletedItem
}
