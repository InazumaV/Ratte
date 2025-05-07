//go:build wireinject

package boot

import (
	"github.com/InazumaV/Ratte/conf"
	"github.com/google/wire"
)

var preSet = wire.NewSet(
	wire.FieldsOf(new(*conf.Conf), "ACME", "Plugin"),
	wire.FieldsOf(new(conf.Plugins), "Core", "Panel"),
	initAcme,
	initCores,
	initPanels,
)
var mainSet = wire.NewSet(
	wire.FieldsOf(new(*conf.Conf), "Node"),
	initNode,
)

type Boot struct {
	Acmes  AcmeGroup
	Cores  CoreGroup
	Panels PanelGroup
	Node   *NodeGroup
}

func (b *Boot) Start() error {
	if err := b.Node.Start(); err != nil {
		return err
	}
	return nil
}

func (b *Boot) Close() error {
	b.Node.Close()
	b.Cores.Close()
	b.Panels.Close()
	return nil
}

func InitBoot(c *conf.Conf) (*Boot, error) {
	wire.Build(
		preSet,
		mainSet,
		wire.Struct(new(Boot), "*"),
	)
	return nil, nil
}
