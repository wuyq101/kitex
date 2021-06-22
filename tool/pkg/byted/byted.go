// Package byted is an extension of KiteX tool for bytedance internal use
package byted

import (
	"path/filepath"

	"github.com/cloudwego/kitex/tool/pkg/generator"
)

var (
	// WithByted is generator feature name
	WithByted = "WithByted"
)

func init() {
	generator.RegisterFeature(WithByted)
	generator.AddGlobalDependency("byted", filepath.Join(generator.KitexImportPath, "byted"))
	generator.AddGlobalMiddleware(AppendBytedImports)
}

// AppendBytedImports add byted imports to generator package information
func AppendBytedImports(next generator.HandleFunc) generator.HandleFunc {
	return func(task *generator.Task, pkg *generator.PackageInfo) (*generator.File, error) {
		switch task.Name {
		case generator.ClientFileName:
			if generator.HasFeature(pkg.Features, WithByted) {
				pkg.AddImports("byted")
			}
			fallthrough
		case generator.ServerFileName, generator.InvokerFileName:
			if generator.HasFeature(pkg.Features, WithByted) {
				pkg.AddImports("byted")
			}
		}
		return next(task, pkg)
	}
}
