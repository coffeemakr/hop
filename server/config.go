package server

import "strconv"

type ListenConfig struct {
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
}

func (c *ListenConfig) GetServerAddress() string {
	addr := c.Host + ":" + strconv.FormatInt(int64(c.Port), 10)
	return addr
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

type AuthenticationConfig struct {
	Key string `json:"key" yaml:"key,omitempty"`
}

type Configuration struct {
	Listen   *ListenConfig         `json:"listen,omitempty" yaml:",omitempty"`
	Database *DatabaseConfig       `json:"database,omitempty" yaml:",omitempty"`
	Auth     *AuthenticationConfig `json:"auth,omitempty" yaml:",omitempty"`
}
