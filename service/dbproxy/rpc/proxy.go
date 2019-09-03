package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	"filestore-server/service/dbproxy/mapper"
	"filestore-server/service/dbproxy/orm"
	dbProxy "filestore-server/service/dbproxy/proto"
)

// DBProxy : DBProxy结构体
type DBProxy struct{}

// ExecuteAction : 请求执行sql函数
func (db *DBProxy) ExecuteAction(ctx context.Context, req *dbProxy.ReqExec, res *dbProxy.RespExec) error {
	resList := make([]orm.ExecResult, len(req.Action))

	// TODO: 检查	req.Sequence req.Transaction两个参数，执行不同的流程
	for idx, singleAction := range req.Action {
		params := []interface{}{}
		dec := json.NewDecoder(bytes.NewReader(singleAction.Params))
		dec.UseNumber()
		// 避免int/int32/int64等自动转换为float64
		//if err := json.Unmarshal(singleAction.Params, &params); err != nil {
		if err := dec.Decode(&params); err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "请求参数有误",
			}
			continue
		}

		for k, v := range params {
			if _, ok := v.(json.Number); ok {
				params[k], _ = v.(json.Number).Int64()
			}
		}

		// 默认串行执行sql函数
		execRes, err := mapper.FuncCall(singleAction.Name, params...)
		if err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "函数调用有误",
			}
			continue
		}
		resList[idx] = execRes[0].Interface().(orm.ExecResult)
	}

	// TODO: 处理异常
	res.Data, _ = json.Marshal(resList)
	return nil
}
