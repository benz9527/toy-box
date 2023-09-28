package cli

import "sync"

type CliKubeResourceMapsDTO struct {
	AddResourcesMap  *sync.Map // create
	ModResourcesMap  *sync.Map // update
	GetResourcesMap  *sync.Map // get
	DelResourcesMap  *sync.Map // delete
	ListResourcesMap *sync.Map // list
}
