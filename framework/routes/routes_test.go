package routes_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestCustomRoutes(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() bool {
			return true
		}
	`
	// root custom routes file
	td.Files["routes/routes.go"] = `
		package routes
		import (
			"github.com/livebud/bud/package/router"
			"app.com/bud/package/controller"
		)
		func Routes(router *router.Router, controller *controller.Controller) {
			router.Get("/custom-route", controller.Index)
		}
	`
	// nested custom routes file (prefixed)
	td.Files["routes/my-prefix/routes.go"] = `
		package routes
		import (
			"github.com/livebud/bud/package/router"
			"app.com/bud/package/controller"
		)
		func Routes(router *router.Router, controller *controller.Controller) {
			router.Get("/nested-custom-route", controller.Index)
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// test custom route
	res, err := app.Get("/custom-route")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: application/json
	`))
	// ensure that this was not mounted at the root
	res, err = app.Get("/nested-custom-route")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff
	`))
	// test custom route with automatic prefix
	res, err = app.Get("/my-prefix/nested-custom-route")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: application/json
	`))
	is.NoErr(app.Close())
}
