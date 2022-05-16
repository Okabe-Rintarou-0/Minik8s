package image

import "github.com/docker/docker/api/types"

type PullConfig struct {
	Verbose bool
	All     bool
}

type ListConfig struct {
	All bool
}
type RemoveConfig = types.ImageRemoveOptions

type DeletedItem = types.ImageDeleteResponseItem

type RemoveResponse struct {
	DeletedItems []DeletedItem
}
