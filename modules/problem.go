package modules

import (
	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/dal"
	"github.com/HeRaNO/cdoj-vjudge/util"
	"github.com/kataras/iris/v12"
)

func SearchProblems(ctx iris.Context) {
	problemName := ctx.URLParamDefault("name", "")
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pageSize", 20)
	if pageSize > 100 {
		ctx.JSON(util.ErrorResponse(config.ErrWrongInfo, "page size too large"))
		return
	}
	if pageSize <= 0 {
		ctx.JSON(util.ErrorResponse(config.ErrWrongInfo, "page size should greater than 0"))
		return
	}
	info, err := dal.GetProblemInfo(ctx, &problemName, (page-1)*pageSize, pageSize)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	ctx.JSON(util.SuccessResponseWithTotal(info, len(info)))
}

func GetProblem(ctx iris.Context) {
	problemID := ctx.URLParamInt64Default("id", 0)
	statement, err := dal.GetStatementByProblemID(ctx, problemID)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	ctx.JSON(util.SuccessResponse(statement))
}
