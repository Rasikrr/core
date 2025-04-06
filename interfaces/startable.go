package interfaces

import (
	"context"
)

type Starter interface {
	Start(ctx context.Context) error
}
