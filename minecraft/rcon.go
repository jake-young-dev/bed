package minecraft

import (
	"fmt"
	"time"

	"github.com/jake-young-dev/mcr"
)

const (
	fullPath = "/data/%s-backup.tar.gz"
	worldDir = "/data/world/"
	filename = "%s-backup.tar.gz"

	//command strings
	saveWorldCommand   = "save-all"
	autosaveOnCommand  = "save-on"
	autosaveOffCommand = "save-off"
	restartCommand     = "restart"
)

/*
A simple wrapper for Minecraft RCon alerts and world backups
*/

type Rcon struct {
	client *mcr.Client
}

type IRcon interface {
	Connect(pwd string) error
	Close() error
	AlertPlayers(alert string) error
	WorldSave() error
	EnableAutosaves() error
	DisableAutosaves() error
	RestartServer() error
}

func NewRconHandler(url string) *Rcon {
	return &Rcon{
		client: mcr.NewClient(url, mcr.WithTimeout(5*time.Second)),
	}
}

func (r *Rcon) Connect(pwd string) error {
	return r.client.Connect(pwd)
}

func (r *Rcon) Close() error {
	return r.client.Close()
}

func (r *Rcon) AlertPlayers(alert string) error {
	_, err := r.client.Command(fmt.Sprintf("say %s", alert))
	return err
}

func (r *Rcon) WorldSave() error {
	_, err := r.client.Command(saveWorldCommand)
	return err
}

func (r *Rcon) EnableAutosaves() error {
	_, err := r.client.Command(autosaveOnCommand)
	return err
}

func (r *Rcon) DisableAutosaves() error {
	_, err := r.client.Command(autosaveOffCommand)
	return err
}

func (r *Rcon) RestartServer() error {
	_, err := r.client.Command(restartCommand)
	return err
}
