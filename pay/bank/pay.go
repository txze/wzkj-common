package bank

import "github.com/mitchellh/mapstructure"

type Pay struct {
	config ConfigBank
}

func (p *Pay) SetConfig(config ConfigBank) {
	p.config = config
}

func (p *Pay) Process(params map[string]string) (*PayResponse, error) {
	params["mchntId"] = p.config.MchntID
	params["notifyUrl"] = p.config.Common.NotifyURL
	params["returnUrl"] = p.config.Common.SyncReturnURL
	signature, err := p.sign(params)
	if err != nil {
		return nil, err
	}

	var rsp PayResponse
	params["signature"] = signature
	_ = mapstructure.Decode(params, &rsp)
	return &rsp, nil
}
