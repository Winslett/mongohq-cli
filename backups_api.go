package main

import (
	"encoding/json"
)

type Backup struct {
	Id             string   `json:"id"`
	CreatedAt      string   `json:"created_at"`
	DatabaseNames  []string `json:"database_names"`
	Status         string   `json:"status"`
	DeploymentSlug string   `json:"deployment"`
	Type           string   `json:"type"`
	Filename       string   `json:"filename"`
	Size           float64  `json:"size"`
	Links          []Hateos `json:"links"`
	Api            Api
}

func (b *Backup) DownloadLink() string {
	for _, link := range b.Links {
		if link.Rel == "download" {
			return link.Href
		}
	}
	return "<Unknown download link>"
}

func (b *Backup) PrettySize() string {
	return prettySize(b.Size)
}

func (api *Api) GetBackups() ([]Backup, error) {
	body, err := api.restGet(api.apiUrl("/accounts/" + api.Config.AccountSlug + "/backups"))

	if err != nil {
		return []Backup{}, err
	}
	var databaseBackupSlice []Backup
	err = json.Unmarshal(body, &databaseBackupSlice)
	return databaseBackupSlice, err
}

func (api *Api) GetBackupsForDeployment(deploymentSlug string) ([]Backup, error) {
	body, err := api.restGet(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deploymentSlug + "/backups"))

	if err != nil {
		return []Backup{}, err
	}
	var databaseBackupSlice []Backup
	err = json.Unmarshal(body, &databaseBackupSlice)
	return databaseBackupSlice, err
}

func (api *Api) GetBackup(backupSlug string) (Backup, error) {
	body, err := api.restGet(api.apiUrl("/accounts/" + api.Config.AccountSlug + "/backups/" + backupSlug))

	if err != nil {
		return Backup{}, err
	}
	var backup Backup
	err = json.Unmarshal(body, &backup)
	return backup, err
}

func (api *Api) RestoreBackup(backup Backup, deploymentName, source, destination string) (Deployment, error) {
	type RestoreBackupParams struct {
		Name           string `json:"name"`
		DatabaseName   string `json:"database_name"`
		SourceDatabase string `json:"source_database"`
	}

	restoreParams := RestoreBackupParams{Name: deploymentName, DatabaseName: destination, SourceDatabase: source}
	data, err := json.Marshal(restoreParams)
	if err != nil {
		return Deployment{}, err
	}

	body, err := api.restPost(api.apiUrl("/accounts/"+api.Config.AccountSlug+"/backups/"+backup.Id+"/restore"), data)
	if err != nil {
		return Deployment{}, err
	}

	var deployment Deployment
	err = json.Unmarshal(body, &deployment)
	return deployment, err
}
