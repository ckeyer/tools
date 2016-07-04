package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	version string
)

type hashOption struct {
	EntireDir     bool
	Path          string
	ToUpper       bool
	HumanReadable bool
	OutputFile    string
	OutputFormat  string
	HMac          string
	Excludes      cli.StringSlice
	Method        string
}

func main() {
	option := &hashOption{ // set default option.
		EntireDir:     false,
		Path:          "",
		ToUpper:       false,
		HumanReadable: false,
		OutputFile:    "",
		OutputFormat:  "{{.FileName}},{{.Size}},{{.Hash}}",
		HMac:          "",
		Excludes:      cli.StringSlice{},
	}

	app := cli.NewApp()
	app.Name = "hash"
	app.Version = version
	app.Usage = "hash command [OPTIONS] filename"
	app.Commands = cli.Commands{}
	for _, hashMethod := range AllHashMethods {
		com := cli.Command{
			Name:  hashMethod,
			Usage: fmt.Sprintf("%s %s [OPTIONS] filename", app.Name, hashMethod),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "R",
					Usage:       "hash entire subtree connected at that point.",
					Destination: &option.EntireDir,
				},
				cli.BoolFlag{
					Name:        "U, toUpper",
					Usage:       "output upper(default is lower)",
					Destination: &option.ToUpper,
				},
				cli.BoolFlag{
					Name:        "H,human",
					Usage:       "Human-readable",
					Destination: &option.HumanReadable,
				},
				cli.StringFlag{
					Name:        "o, output",
					Usage:       "output file",
					Destination: &option.OutputFile,
				},
				cli.StringFlag{
					Name:        "f, format",
					Usage:       "output format,.",
					Destination: &option.OutputFormat,
				},
				cli.StringFlag{
					Name:        "hmac",
					Usage:       "hmachash",
					Destination: &option.HMac,
				},
				cli.StringSliceFlag{
					Name:  "E,exclude",
					Usage: "exclude",
					Value: &option.Excludes,
				},
			},
			Before: func(c *cli.Context) error { // check input args
				if len(c.Args()) != 1 {
					fmt.Println("error")
					return fmt.Errorf("invalid filename")
				}
				path := c.Args().Get(0)
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				option.Path = absPath

				return nil
			},
			Action: func(c *cli.Context) error {
				option.Method = fmt.Sprintf("%#v", c.Command.Name)

				return Run(option)
			},
		}
		app.Commands = append(app.Commands, com)
	}

	app.Run(os.Args)
}

func Run(option *hashOption) error {

	if !option.EntireDir {

	}

	fmt.Printf("debug option: %#v\n", option)
	return nil
}
