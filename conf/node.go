package conf

import (
	"fmt"
	"github.com/goccy/go-json"
	"strconv"
	"strings"
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
	Core   string                 `json:"CorePlugin"`
	Panel  string                 `json:"Panel"`
	Acme   string                 `json:"Acme"`
	Cert   Cert                   `json:"Cert"`
	Hook   Hook                   `json:"Hook"`
	Limit  Limit                  `json:"Limit"`
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

type Limit struct {
	IPLimit    int      `json:"IPLimit"`
	SpeedLimit IntBytes `json:"SpeedLimit"`
}

type IntBytes uint64

func (b *IntBytes) UnmarshalJSON(data []byte) error {
	var num uint64
	err := json.Unmarshal(data, &num)
	if err == nil {
		*b = IntBytes(num)
	}
	var numS string
	err = json.Unmarshal(data, &numS)
	if err != nil {
		return fmt.Errorf("failed to unmarshal intBytes: %v", err)
	}
	unit := numS[len(numS)-2:]
	num, err = strconv.ParseUint(numS[:len(numS)-2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid bytes num: %s", numS)
	}
	switch strings.ToLower(unit) {
	case "kb":
		*b = IntBytes(num) * 1000 / 8
	case "mb":
		*b = IntBytes(num) * 1000 * 1000 / 8
	case "gb":
		*b = IntBytes(num) * 1000 * 1000 * 1000 / 8
	case "tb":
		*b = IntBytes(num) * 1000 * 1000 * 1000 * 1000 * 1000 / 8
	default:
		return fmt.Errorf("invalid bytes unit: %s", unit)
	}
	*b = IntBytes(num)
	return nil
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
		return fmt.Errorf("failed to unmarshal Node: %v", err)
	}

	n.Remote = Remote{
		APIHost: "http://127.0.0.1",
		Timeout: 30,
	}
	if len(rn.RemoteRaw) > 0 {
		err = json.Unmarshal(rn.RemoteRaw, &n.Remote)
		if err != nil {
			return fmt.Errorf("failed to unmarshal RemoteRaw: %v", err)
		}
	} else {
		err = json.Unmarshal(data, &n.Remote)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Remote: %v", err)
		}
	}
	n.Options = Options{}
	if len(rn.OptRaw) > 0 {
		err = json.Unmarshal(rn.OptRaw, &n.Options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal OptRaw: %v", err)
		}
	} else {
		err = json.Unmarshal(data, &n.Options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Options: %v", err)
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
			return fmt.Errorf("failed to unmarshal TriggerRaw: %v", err)
		}
	} else {
		err = json.Unmarshal(data, &n.Trigger)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Trigger: %v", err)
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
