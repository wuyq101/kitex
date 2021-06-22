package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex/tool/pkg/pluginmode/protoc"
	"github.com/cloudwego/kitex/tool/pkg/pluginmode/thriftgo"
)

func main() {
	BuildKiteX()
}

func BuildKiteX() {
	// run as a plugin
	switch filepath.Base(os.Args[0]) {
	case thriftgo.PluginName:
		os.Exit(thriftgo.Run())
	case protoc.PluginName:
		os.Exit(protoc.Run())
	}

	// run as kitex
	args.parseArgs()

	out := new(bytes.Buffer)
	cmd := buildCmd(&args, out)
	err := cmd.Run()
	if err != nil {
		if args.Use != "" {
			out := strings.TrimSpace(out.String())
			if strings.HasSuffix(out, thriftgo.TheUseOptionMessage) {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}
