package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/hello",
			Handler: func(c echo.Context) error {
				return c.String(200, "Hello world!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	app.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		if e.Record.Collection().Name == "drafts" {
			// app.Dao().DB().Update("drafts", dbx.Params{"is_publish": false}, dbx.And(dbx.HashExp{"id_user": e.Record.GetDataValue("id_user")}, dbx.NewExp(fmt.Sprintf("id != %s", e.Record.Id)))).Execute()
			log.Println("UPDATE drafts SET is_publish = 0 WHERE is_publish = 1 AND id_user = " + fmt.Sprintf("'%s'", e.Record.Get("id_user")))
			query := app.DB().NewQuery("UPDATE drafts SET is_publish = 0 WHERE is_publish = 1 AND id_user = " + fmt.Sprintf("'%s'", e.Record.Get("id_user")))
			if _, err := query.Execute(); err != nil {
				return err
			}
		}
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
