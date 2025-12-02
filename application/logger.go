package application

import (
	"github.com/Rasikrr/core/log"
)

func (a *App) InitLogger() error {
	return log.Init(
		a.Config().Env(),
		a.Config().Logger.Level,
		a.Config().Logger.AddSource,
	)
}
