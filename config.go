package main

import (
  "github.com/BurntSushi/toml"
)

type mysqlconfig struct {
  Enabled bool
  Connectstring string
}

type redisconfig struct {
  Enabled bool
  Host string
  Password string
}

type serverconfig struct {
  Maindb string
  Ip string
  Port int
}

type config struct {
  Mysqlconfig mysqlconfig
  Redisconfig redisconfig
  Serverconfig serverconfig
}

func loadConfig(file string) (*config, error) {
  newconfig := new(config)
  if _, err := toml.DecodeFile(file, &newconfig); err != nil {
    return newconfig, err
  }
  return newconfig, nil // config, no error
}
