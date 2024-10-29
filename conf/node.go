package conf

import (
	"fmt"
	"github.com/goccy/go-json"
)

type rawNodeConfig struct {
	Name       string          `json:"Name"`
	RemoteRaw  json.RawMessage `json:"Remote"`
	OptRaw     json.RawMessage `json:"Options"`
	TriggerRaw json.RawMessage `json:"Trigger"`
}

type Remote struct {
	APIHost  string `json:"ApiHost"`
	NodeID   int    `json:"NodeID"`
	Key      string `json:"ApiKey"`
	NodeType string `json:"NodeType"`
	Timeout  int    `json:"Timeout"`
}

type Options struct {
	Core   string                 `json:"Core"`
	Panel  string                 `json:"Panel"`
	Acme   string                 `json:"Acme"`
	Cert   Cert                   `json:"Cert"`
	Hook   Hook                   `json:"Hook"`
	Expand map[string]interface{} `json:"Other"`
}

type Trigger struct {
	PullNodeCron   any `json:"PullNodeCron"`
	PullUserCron   any `json:"PullUserCron"`
	ReportUserCron any `json:"ReportUserCron"`
	RenewCertCron  any `json:"RenewCertCron"`
}

type Cert struct {
	Domain   string `json:"Domain"`
	CertPath string `json:"Cert"`
	KeyPath  string `json:"Key"`
}

type Hook struct {
	BeforeAddNode string `json:"BeforeAddNode"`
	AfterAddNode  string `json:"AfterAddNode"`
	BeforeDelNode string `json:"BeforeDelNode"`
	AfterDelNode  string `json:"AfterDelNode"`
}

type Node struct {
	Name    string  `json:"Name"`
	Remote  Remote  `json:"-"`
	Trigger Trigger `json:"-"`
	Options Options `json:"-"`
}

func (n *Node) UnmarshalJSON(data []byte) (err error) {
	rn := rawNodeConfig{}
	err = json.Unmarshal(data, &rn)
	if err != nil {
		return err
	}

	n.Remote = Remote{
		APIHost: "http://127.0.0.1",
		Timeout: 30,
	}
	if len(rn.RemoteRaw) > 0 {
		err = json.Unmarshal(rn.RemoteRaw, &n.Remote)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(data, &n.Remote)
		if err != nil {
			return
		}
	}
	n.Options = Options{}
	if len(rn.OptRaw) > 0 {
		err = json.Unmarshal(rn.OptRaw, &n.Options)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(data, &n.Options)
		if err != nil {
			return
		}
	}
	n.Trigger = Trigger{
		PullNodeCron:   60,
		PullUserCron:   60,
		ReportUserCron: 60,
		RenewCertCron:  "0 2 * * *",
	}
	if len(rn.TriggerRaw) > 0 {
		err = json.Unmarshal(rn.OptRaw, &n.Trigger)
		if err != nil {
			return
		}
	} else {
		err = json.Unmarshal(data, &n.Trigger)
		if err != nil {
			return
		}
	}
	if len(rn.Name) > 0 {
		n.Name = rn.Name
	} else {
		n.Name = fmt.Sprintf("{T:%s;A:%s;I:%d;}",
			n.Remote.NodeType,
			n.Remote.APIHost,
			n.Remote.NodeID)
	}
	return
}
