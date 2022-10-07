package log

import (
	"go.uber.org/zap"

	def "StakeBackendGoTest/pkg/const"
)

var Logger *zap.Logger

func Init(env def.EnvMode, options ...zap.Option) error {
	var err error
	if env == def.ProdMode {
		Logger, err = zap.NewProduction(options...)
	} else {
		Logger, err = zap.NewDevelopment(options...)
	}

	return err
}
