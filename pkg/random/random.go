package random

import (
	"github.com/google/uuid"
)

func UUID4String() string {
	return uuid.New().String()
}
