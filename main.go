package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

var (
	appName, appVer string
)

func main() {

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "CLI tools which used to visualize Akamai Network list usage."
	app.Version = appVer
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Petr Artamonov",
		},
		{
			Name: "Rafal Pieniazek",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Value: "json",
			Usage: "Output result as .json metadata or as .dot graphviz files",
		},
		cli.StringFlag{
			Name:  "map-file, m",
			Value: "",
			Usage: "Path for `configuration map file`",
		},
		cli.StringFlag{
			Name:  "id, i",
			Value: "",
			Usage: "Network List `ID`",
		},
		cli.StringFlag{
			Name:  "name, n",
			Value: "",
			Usage: "Network List `NAME`",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "",
			Usage: "Path for destination files. If empty results will be sent to STDOUT",
		},
		cli.StringFlag{
			Name:  "source, s",
			Value: os.TempDir(),
			Usage: "Path, where Security configuration JSON files are located",
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	app.Action = cmdSearch

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}

}
