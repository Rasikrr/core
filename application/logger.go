package application

import (
	"github.com/Rasikrr/core/log"
)

func (a *App) InitLogger() {
	log.Init(a.Config().Env())
}
