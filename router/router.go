package router

import (
	"github.com/HeRaNO/cdoj-vjudge/dal"
	"github.com/HeRaNO/cdoj-vjudge/modules"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/basicauth"
)

func InitApp() *iris.Application {
	app := iris.New()
	allowFunc := func(ctx iris.Context, username string, password string) (interface{}, bool) {
		ok, err := dal.IsAuthValid(ctx, &username, &password)
		if err != nil {
			return nil, false
		}
		return nil, ok
	}
	opt := basicauth.Options{
		Realm:        basicauth.DefaultRealm,
		Allow:        allowFunc,
		ErrorHandler: basicauth.DefaultErrorHandler,
	}

	app.Use(basicauth.New(opt))
	app.Get("/searchProblems", modules.SearchProblems)
	app.Get("/getProblem", modules.GetProblem)
	app.Post("/submit", modules.Submit)
	app.Get("/status", modules.GetStatus)
	return app
}
