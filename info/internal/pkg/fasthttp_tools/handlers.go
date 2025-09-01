package fasthttp_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"info/internal/pkg/apperror"
	"strconv"
)

const (
	AuthClientKey = "http.client"
	RequestIdKey  = "http.requestId"
)

var stopList = map[string]struct{}{}

func ParseQueryArgUint(ctx *fasthttp.RequestCtx, name string) (uint, error) {
	valStr, err := ParseQueryArgString(ctx, name)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseUint(valStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("[%w] failed to parse uint param %s; error: %w", apperror.ErrBadRequest, name, err)
	}

	return uint(val), nil
}

func ParseQueryArgString(ctx *fasthttp.RequestCtx, name string) (string, error) {
	val := string(ctx.QueryArgs().Peek(name))
	if val == "" {
		return "", fmt.Errorf("[%w] empty param %s", apperror.ErrNotFound, name)
	}

	return val, nil
}

func BadRequest(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	resp := errorResp{
		Error:     err.Error(),
		RequestId: "",
	}
	if id, ok := ctx.UserValue(RequestIdKey).(string); ok {
		resp.RequestId = id
	}

	res, _ := json.Marshal(resp)

	ctx.Write(res)
}

func InternalError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	resp := errorResp{
		Error:     err.Error(),
		RequestId: "",
	}
	if id, ok := ctx.UserValue(RequestIdKey).(string); ok {
		resp.RequestId = id
	}

	res, _ := json.Marshal(resp)
	ctx.Write(res)
}

type errorResp struct {
	Error     string
	RequestId string
}

func Success(ctx *fasthttp.RequestCtx, body []byte) error {
	if body == nil {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil
	}
	if _, err := ctx.Write(body); err != nil {
		return fmt.Errorf("Write(body) error: %w", err)
	}
	return nil
}

func parseUserId(ctx *fasthttp.RequestCtx) (string, error) {
	v, ok := ctx.UserValue("user_id").(string)
	if !ok {
		return "", errors.New("invalid userId")
	}

	if len(v) > 20 {
		return "", errors.New("invalid userId")
	}

	return v, nil
}
