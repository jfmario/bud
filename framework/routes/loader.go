
package routes

import (

	"io/fs"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module, parser *parser.Parser) (*State, error) {
	return &State{}, nil
}