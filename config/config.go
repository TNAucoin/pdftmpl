package config

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type Conf struct {
	Server ConfServer
	Goten  ConfGoten
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type ConfGoten struct {
	Host          string `env:"GOTEN_HOST,required"`
	Port          int    `env:"GOTEN_PORT,required"`
	VolumeOutPath string `env:"OUT_VOLUME_PATH,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	return &c
}
