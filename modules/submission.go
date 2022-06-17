package modules

import (
	"context"
	"log"

	workconfig "github.com/HeRaNO/cdoj-execution-worker/config"
	workmodel "github.com/HeRaNO/cdoj-execution-worker/model"
	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/dal"
	"github.com/HeRaNO/cdoj-vjudge/model"
	"github.com/HeRaNO/cdoj-vjudge/util"
	"github.com/bytedance/sonic"
	"github.com/kataras/iris/v12"
)

func Submit(ctx iris.Context) {
	req := model.Submission{}
	if err := ctx.UnmarshalBody(&req, iris.UnmarshalerFunc(sonic.Unmarshal)); err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	checker, limit, err := dal.GetProblemCheckerInfo(ctx, req.ProblemID)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	execReq, err := util.MakeSubmission(req, checker, limit)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrWrongInfo, err.Error()))
		return
	}
	id, err := dal.SendMessage(execReq)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	ctx.JSON(util.SuccessResponse(id))
}

func GetStatus(ctx iris.Context) {
	id := ctx.URLParamDefault("submission_id", "")
	status, err := dal.GetSubmissionResult(ctx, &id)
	if err != nil {
		ctx.JSON(util.ErrorResponse(config.ErrInternal, err.Error()))
		return
	}
	ctx.JSON(util.SuccessResponse(status))
}

func UpdateStatus(resp workmodel.Response, corId string) {
	res := model.Result{}
	switch resp.ErrCode {
	case workconfig.CE:
		res = util.MakeCEResult(resp.Data)
	case workconfig.RE:
		res = util.MakeREResult(resp.ErrMsg, resp.Data)
	case workconfig.IE:
		res = util.MakeIEResult()
	case workconfig.OK:
		if resp.ErrMsg != "running" && resp.ErrMsg != "success" && resp.ErrMsg != "wrong answer" {
			res = util.MakeUnknownResult("unknown error message")
		} else {
			res = util.MakeOKResult(resp.ErrMsg, resp.Data)
		}
	default:
		res = util.MakeUnknownResult("unknown error code")
	}
	resS, err := sonic.MarshalString(res)
	if err != nil {
		log.Printf("[ERROR] UpdateStatus(): marshal error: %s", err.Error())
		return
	}
	dal.SetSubmissionResult(context.Background(), &corId, &resS)
}

func ListenMQ() {
	for resp := range config.MQ {
		res := workmodel.Response{}
		err := sonic.Unmarshal(resp.Body, &res)
		resp.Ack(false)
		log.Printf("[INFO] mq received msg: %+v", res)
		if err != nil {
			log.Printf("[ERROR] ListenMQ(): unmarshal error: %s", err.Error())
			continue
		}
		UpdateStatus(res, resp.CorrelationId)
	}
}
