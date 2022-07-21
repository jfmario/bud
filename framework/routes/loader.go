package routes

import (
	"io/fs"
	"path"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/vfs"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module, parser *parser.Parser) (*State, error) {
	exist, err := vfs.SomeExist(fsys, "routes")
	if err != nil {
		return nil, err
	} else if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	loader := &loader{
		fsys:   fsys,
		module: module,
		parser: parser,
	}
	return loader.Load()
}

// loader struct
type loader struct {
	bail.Struct
	fsys   fs.FS
	module *gomod.Module
	parser *parser.Parser
}

func (l *loader) Load() (state *State, err error) {

	defer l.Recover2(&err, "routes: unable to load")

	state = new(State)
	state.ControllerImport = l.module.Import("bud/package/controller")
	state.RoutesFuncs = l.loadRoutesFuncs("routes", "", make([]*RoutesFunc, 0))
	if len(state.RoutesFuncs) > 0 {
		state.HasRoutes = true
	}
	return state, nil
}

// load route functions from routes
func (l *loader) loadRoutesFuncs(routesPath, prefix string, routesFuncs []*RoutesFunc) []*RoutesFunc {

	des, err := fs.ReadDir(l.fsys, routesPath)
	if err != nil {
		l.Bail(err)
	} else if len(des) == 0 {
		l.Bail(fs.ErrNotExist)
	}

	// routesFunc := new(RoutesFunc)

	shouldParse := false

	for _, de := range des {
		if !de.IsDir() && valid.RoutesFile(de.Name()) {
			shouldParse = true
			continue
		}
		if de.IsDir() && valid.Dir(de.Name()) {
			routesFuncs = l.loadRoutesFuncs(path.Join(routesPath, de.Name()), prefix+"/"+de.Name(), routesFuncs)
		}
	}

	if !shouldParse {
		return routesFuncs
	}

	pkg, err := l.parser.Parse(routesPath)
	if err != nil {
		l.Bail(err)
	}

	// look for function called Routes
	rawRoutesFunc := pkg.PublicFunction("Routes")
	if rawRoutesFunc == nil {
		return routesFuncs
	}

	// must have 2 params
	rawRoutesFuncParams := rawRoutesFunc.Params()
	if len(rawRoutesFuncParams) != 2 {
		return routesFuncs
	}

	// types must be *router.Router and *controller.Controller
	if rawRoutesFuncParams[0].Type().String() != "*router.Router" || rawRoutesFuncParams[1].Type().String() != "*controller.Controller" {
		return routesFuncs
	}

	importPath, err := pkg.Import()
	if err != nil {
		l.Bail(err)
	}

	routesFuncs = append(routesFuncs, &RoutesFunc{
		Import: importPath,
		Prefix: prefix,
	})

	return routesFuncs
}
