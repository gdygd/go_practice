package cmd

import (
	"context"
	"fmt"
	"go_redis/config"
	"go_redis/general/cache"
)

type CLI struct {
	ctx       context.Context
	cmd       []string
	command   string
	args      []string
	Exit      bool // cli exit
	Terminate bool // process terminate
	Rdb       *cache.RedisClient
	Config    *config.Config
}

func NewCLI() *CLI {
	config, err := initConfig()
	if err != nil {
		return nil
	}

	rdb := cache.NewRedisClient(config.RedisAddr)
	return &CLI{
		ctx:     context.Background(),
		cmd:     make([]string, 0, 0),
		command: "",
		Rdb:     rdb,
		Config:  &config,
	}
}

func initConfig() (config.Config, error) {

	return config.LoadConfig(".")
}

func (c *CLI) SetCommand(commands []string) {
	c.cmd = make([]string, 0, 0)
	for _, val := range commands {
		c.cmd = append(c.cmd, val)
	}
}

func (c *CLI) PrintCmd() {
	fmt.Printf("commands : %v\n", c.cmd)
}

func (c *CLI) Run() {
	c.command = c.cmd[0]
	c.args = make([]string, len(c.cmd)-1)
	if len(c.cmd) > 1 {
		c.args = c.cmd[1:]
	}

	c.check()
}

func (c *CLI) check() {
	if c.command == "help" {
		c.help()
	} else if c.command == "system" {
		c.system()
	} else if c.command == "version" {
		c.version()
	} else if c.command == "process" {
		c.process()
	} else if c.command == "debug" {
		c.debug()
	} else if c.command == "terminate" {
		c.terminate()
	} else {
		fmt.Println("Invlid command")
	}
}
