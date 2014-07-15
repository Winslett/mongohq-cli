package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func (c *Controller) ListBackups(deploymentName string) {
	backupsSlice, err := c.Api.GetBackups(deploymentName)

	if err != nil {
		fmt.Println("Error retreiving backups: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("== Backups")
	for _, backup := range backupsSlice {
		fmt.Println(backup.Filename)
	}
}

func (c *Controller) ShowBackup(filename string) {
	backup, err := c.findBackupByFilename(filename)
	if err != nil {
		fmt.Println("Error retreiving backup: " + err.Error())
		os.Exit(1)
	}
	deployment, _ := c.Api.GetDeployment(backup.DeploymentId)
	fmt.Println("== Backup " + filename)
	fmt.Println(" deployment    : " + deployment.Name)
	fmt.Println(" databases     : " + strings.Join(backup.DatabaseNames, ", "))
	fmt.Println(" created at    : " + backup.CreatedAt)
	fmt.Println(" type          : " + backup.Type)
	fmt.Println(" size          : " + backup.PrettySize())
	fmt.Println(" download      : " + backup.DownloadLink())
}

func (c *Controller) RestoreBackup(filename, source, destination string) {
	backup, err := c.findBackupByFilename(filename)
	if err != nil {
		fmt.Println("Error retreiving backup: " + err.Error())
		os.Exit(1)
	}

	deployment, err := c.Api.RestoreBackup(backup, source, destination)
	if err != nil {
		fmt.Println("Error restoring backup: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("=== Restoring from database " + source + " in backup file " + filename + " to new database " + destination)
	c.pollNewDeployment(deployment)
}

func (c *Controller) findBackupByFilename(filename string) (Backup, error) {
	backups, err := c.Api.GetBackups("<string>") // <string> is a default that will be ignored, thus requesting all backups
	if err != nil {
		return Backup{}, err
	}
	for _, backup := range backups {
		if backup.Filename == filename {
			return backup, err
		}
	}
	return Backup{}, errors.New("Could not find backup with value " + filename)
}
