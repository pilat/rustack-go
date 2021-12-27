package rustack

import "fmt"

type StorageProfile struct {
	manager *Manager
	ID      string `json:"id"`
	Name    string `json:"name"`
}

func (v *Vdc) GetStorageProfiles() (storageProfiles []*StorageProfile, err error) {
	args := Arguments{
		"vdc": v.ID,
	}

	path := "v1/storage_profile"
	err = v.manager.GetItems(path, args, &storageProfiles)
	for i := range storageProfiles {
		storageProfiles[i].manager = v.manager
	}
	return
}

func (v *Vdc) GetStorageProfile(id string) (storageProfile *StorageProfile, err error) {
	args := Arguments{
		"vdc": v.ID,
	}

	path := fmt.Sprintf("v1/storage_profile/%s", id)
	err = v.manager.Get(path, args, &storageProfile)
	storageProfile.manager = v.manager
	return
}
