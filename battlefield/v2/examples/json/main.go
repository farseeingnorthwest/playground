package main

import (
	"encoding/json"
	"os"

	"github.com/alecthomas/kong"
	"github.com/farseeingnorthwest/playground/battlefield/v2/examples"
)

type EffectCmd struct {
	Name string `arg`
}

func (c EffectCmd) Run(encoder *json.Encoder) error {
	return encoder.Encode(examples.Effect[c.Name])
}

type RegularCmd struct {
	Index int `arg`
}

func (c RegularCmd) Run(encoder *json.Encoder) error {
	return encoder.Encode(examples.Regular[c.Index])
}

type SpecialCmd struct {
	Group int `arg`
	Index int `arg`
}

func (c SpecialCmd) Run(encoder *json.Encoder) error {
	return encoder.Encode(examples.Special[c.Group][c.Index])
}

func main() {
	var cli struct {
		Effect  EffectCmd  `cmd`
		Regular RegularCmd `cmd`
		Special SpecialCmd `cmd`
		Indent  bool       `short:"i"`
	}

	ctx := kong.Parse(&cli)
	encoder := json.NewEncoder(os.Stdout)
	if cli.Indent {
		encoder.SetIndent("", "  ")
	}

	ctx.FatalIfErrorf(ctx.Run(encoder))
}
