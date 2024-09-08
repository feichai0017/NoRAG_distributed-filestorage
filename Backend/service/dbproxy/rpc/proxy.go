package rpc

import (
	"bytes"
	"cloud_distributed_storage/Backend/service/dbproxy/mapper"
	"cloud_distributed_storage/Backend/service/dbproxy/orm"
	dbproxy "cloud_distributed_storage/Backend/service/dbproxy/proto"
	"context"
	"encoding/json"
)

type DBProxy struct{}

func (d *DBProxy) ExecuteAction(ctx context.Context, req *dbproxy.ReqExec, res *dbproxy.ResExec) error {
	resList := make([]orm.ExecResult, len(req.Actions))

	for index, singleAction := range req.Actions {
		var params []interface{}
		dec := json.NewDecoder(bytes.NewReader(singleAction.Params))
		dec.UseNumber()
		if err := dec.Decode(&params); err != nil {
			resList[index] = orm.ExecResult{
				Suc: false,
				Msg: "Invalid params",
			}
			continue
		}

		for k, v := range params {
			if _, ok := v.(json.Number); ok {
				params[k], _ = v.(json.Number).Int64()
			}
		}

		execRes, err := mapper.FunCall(singleAction.Name, params...)
		if err != nil {
			resList[index] = orm.ExecResult{
				Suc: false,
				Msg: "function call failed",
			}
			continue
		}
		resList[index] = execRes[0].Interface().(orm.ExecResult)
	}

	var err error
	res.Data, err = json.Marshal(resList)
	if err != nil {
		res.Msg = "Failed to marshal resList"
		return nil
	}
	return nil
}
