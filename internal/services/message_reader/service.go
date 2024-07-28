package messagereader

import (
	"context"
)

//go:generate mockgen -source=service.go -destination=../mock/servicereader.go
type ServiceReader interface {
	StartConsumerProcessingMessage(ctx context.Context)
}
