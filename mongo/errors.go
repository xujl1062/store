package mongo

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

var (
	InvalidKeyError = errors.New("Invalid key , the key should be db/collection fmt")
)

type InvalidPtrError struct {
	Typ reflect.Type
}

func (err *InvalidPtrError) Error() string {
	return fmt.Sprintf("Entity type should be ptr but received type is %s", err.Typ.String())
}
