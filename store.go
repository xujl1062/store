package store

import "context"

type Store interface {
	Get(ctx context.Context, key string, criteria map[string]interface{}, out interface{}) error
	Save(ctx context.Context, key string, entity interface{}) error
	List(ctx context.Context, key string, criteria map[string]interface{}, target interface{}) ([]interface{}, error)
	Update(ctx context.Context, key string, criteria map[string]interface{}, entity interface{}) error
}

type ExtOptions struct {
	Options []interface{}
}

func (e *ExtOptions) Apply(opt interface{}) *ExtOptions {
	e.Options = append(e.Options, opt)
	return e
}
