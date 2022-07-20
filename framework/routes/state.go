
package routes

type State struct {
	ControllerImport string
	HasRoutes bool
	RoutesFuncs []*RoutesFunc
}

type RoutesFunc struct {
	Import string
	Prefix string
}