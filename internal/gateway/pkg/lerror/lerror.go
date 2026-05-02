package lerror

import (
	"seckill/internal/gateway/dto"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

func GenErrorResponse(err error) dto.Response {
	var resp dto.Response
	if err != nil {
		// 如果是Rpc业务报错
		if bizErr, ok := kerrors.FromBizStatusError(err); ok {
			resp = dto.Response{Status: bizErr.BizStatusCode(), Info: bizErr.BizMessage()}
		} else {
			resp = dto.InternalError(err)
		}
	}
	return resp
}
