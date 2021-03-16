package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "test cmd",
		Action: func(c *cli.Context) error {
			fmt.Println("this is a test cmd")
			return nil
		},
	},
	{
		Name:    "new",
		Aliases: []string{"fuck"},
		Usage:   "创建一个属于你的JoGo 框架",
		Action: func(c *cli.Context) error {
			fmt.Println("你的JoGo框架已经成功生成～")
			return nil
		},
	},

	{
		Name:    "swagInit",
		Aliases: []string{"swag"},
		Usage:   "生成你的swagger doc 文档",
		Action: func(c *cli.Context) error {
			fmt.Println("改项目的文档成功生成～")
			return nil
		},
	},
}
