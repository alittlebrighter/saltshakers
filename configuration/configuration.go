package configuration

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	"github.com/alittlebrighter/saltshakers/messages"
)

func Producer() actor.Actor {
	return NewConfigurationActor()
}

// This is a fairly simple application so for the most part these values are not configurable
var appConfiguration = AppConfig{
	IO: []IOConfig{
		IOConfig{
			kind: Wails,
			GenericConfig: GenericConfig{
				Params: map[string]interface{}{
					"wailsConfig": wails.AppConfig{
						Width:  1024,
						Height: 768,
						Title:  "Saltshakers",
						JS:     mewn.String("../frontend/dist/app.js"),
						CSS:    mewn.String("../frontend/dist/app.css"),
						Colour: "#131313",
					},
				},
			},
		},
	},
	Persistence: []PersistenceConfig{
		PersistenceConfig{
			kind: Bolt,
			GenericConfig: GenericConfig{
				Params: map[string]interface{}{
					"dbPath": "./saltshakers.db",
				},
			},
		},
		PersistenceConfig{
			kind: File,
			GenericConfig: GenericConfig{
				Params: map[string]interface{}{
					"filePath": "%s/saltshakers.log",
				},
			},
		},
	},
}

type ConfigurationActor struct {
	config AppConfig
}

func (state *ConfigurationActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case messages.GetIOConfiguration:
		context.Respond(state.config.IO)
	case messages.GetPersistenceConfiguration:
		context.Respond(state.config.Persistence)
	case *actor.Started:
		state.config = appConfiguration
		// start a file watcher?
	case *actor.Stopping:
	case *actor.Stopped:
	case *actor.Restarting:
	}
}

func NewConfigurationActor() *ConfigurationActor {
	return &ConfigurationActor{}
}

type GenericConfig struct {
	Params map[string]interface{}
}

type IOConfig struct {
	kind IOKind
	GenericConfig
}

func (c *IOConfig) Kind() IOKind {
	return c.kind
}

type PersistenceConfig struct {
	kind PersistenceKind
	GenericConfig
}

func (c *PersistenceConfig) Kind() PersistenceKind {
	return c.kind
}

type AppConfig struct {
	IO          []IOConfig
	Persistence []PersistenceConfig
}

type IOKind uint8

const (
	Wails IOKind = iota
)

type PersistenceKind uint8

const (
	File PersistenceKind = iota
	Bolt
)

type PersistenceUse uint8

const (
	AppData PersistenceUse = iota
	UserConfig
	Cache
)
