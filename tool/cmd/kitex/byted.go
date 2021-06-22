package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex"
	"github.com/cloudwego/kitex/tool/pkg/byted"
	"github.com/cloudwego/kitex/tool/pkg/log"
	"github.com/cloudwego/kitex/tool/pkg/selfupdate"
	"github.com/cloudwego/kitex/tool/pkg/util"
)

func init() {
	var onlyCore, disableSelfUpdate bool
	args.addExtraFlag(&extraFlag{
		apply: func(f *flag.FlagSet) {
			f.BoolVar(&onlyCore, "core", false,
				"Generate codes without injecting the byted suite -- the default service governing functionality set.")

			f.BoolVar(&disableSelfUpdate, "disable-self-update", false,
				"Disable self-update this time.")
		},
		check: func(a *arguments) {
			if !onlyCore {
				if !a.AddFeature(byted.WithByted) {
					log.Warn("add byted feature failed")
				}
			}

			if !disableSelfUpdate {
				selfupdate.Start("kitex", kitex.Version, true)
			}

			tool, version := "thriftgo", "0.0.0"
			if a.IDLType == "protobuf" {
				tool = "protoc"
			}

			path, err := exec.LookPath(tool)
			if err == nil {
				version, err = queryVersion(path)
				if err != nil {
					log.Warnf("Failed to query version of '%s': %s\n", path, err.Error())
					os.Exit(1)
				}
			} else {
				path = filepath.Join(util.GetGOPATH(), "bin", tool)
			}

			version = strings.TrimSpace(version)
			if !strings.HasPrefix(version, "v") {
				version = "v" + version
			}

			if version != "v0.0.0" && disableSelfUpdate {
				return
			}

			updated, version2, err := selfupdate.Update(tool, version, path)
			if err != nil {
				log.Warnf("Failed to update '%s': %s\n", tool, err.Error())
				if version == "v0.0.0" {
					log.Warnf("command %s not found", tool)
					os.Exit(1)
				}
			}
			if updated {
				log.Warnf("Update %s from %s to %s\n", tool, version, version2)
			}
		},
	})
}

func queryVersion(exe string) (version string, err error) {
	var buf strings.Builder
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe, "--version",
		},
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: &buf,
	}
	err = cmd.Run()
	if err == nil {
		version = strings.Split(buf.String(), " ")[1]
	}
	return
}
