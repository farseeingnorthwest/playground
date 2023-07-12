package main

import (
	"encoding/base64"
	"github.com/alecthomas/kong"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
)

func main() {
	var cli struct {
		V1     bool `name:"v1" xor:"version"`
		V2     bool `name:"v2" xor:"version"`
		Base58 bool `name:"base58" xor:"encoding"`
		Base64 bool `name:"base64" xor:"encoding" help:"URL-safe base64 encoding."`
	}
	kong.Parse(&cli, kong.Description("UUID (default v4) generator."))

	var (
		id  uuid.UUID
		err error
	)
	switch {
	case cli.V1:
		id = uuid.New()
	case cli.V2:
		id, err = uuid.NewDCEGroup()
	default:
		id, err = uuid.NewRandom()
	}
	if err != nil {
		panic(err)
	}

	switch {
	case cli.Base58:
		println(base58.Encode(id[:]))
	case cli.Base64:
		println(base64.URLEncoding.EncodeToString(id[:]))
	default:
		println(id.String())
	}
}
