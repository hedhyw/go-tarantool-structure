package getter

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hedhyw/go-tarantool-structure/internal/consts"
	"github.com/hedhyw/go-tarantool-structure/internal/errs"
	"github.com/hedhyw/go-tarantool-structure/internal/interfaces"
	"github.com/mitchellh/mapstructure"
	tarantool "github.com/tarantool/go-tarantool"
)

// InvalidEntityError describes an entity parsing error
type InvalidEntityError struct{}

func (e InvalidEntityError) Error() string {
	return "Invalid entity "
}

// getter implements IGetter
type getter struct {
	s string      // space
	m interface{} // model
	c interfaces.IConnection
}

// IGetter is a helper for records getting
type IGetter interface {
	All(
		index interface{},
		val interface{},
		offset uint32,
		limit uint32,
	) ([]interface{}, error)
	FindBy(index interface{}, val interface{}) (interface{}, error)
	Find(id string) (interface{}, error)
}

// NewGetter wraps a model with a getter methods
func NewGetter(
	space string,
	model interface{},
	conn interfaces.IConnection,
) IGetter {
	return &getter{space, model, conn}
}

// model parses tuple into the model
func (g getter) model(t []interface{}) (interface{}, error) {
	var in = map[string]interface{}{}
	ref := reflect.TypeOf(g.m).Elem()

	// Convert tuple into the map
	for i, count := 0, ref.NumField(); i < count; i++ {
		f := ref.Field(i)
		val, ok := f.Tag.Lookup(consts.TagName)
		if !ok {
			continue
		}
		index, err := strconv.Atoi(val)
		if err != nil {
			return nil, errs.NewInvalidTagValueError(f.Name)
		}
		in[f.Name] = t[index]
	}

	if err := mapstructure.Decode(in, g.m); err != nil {
		return nil, err
	}

	return g.m, nil
}

// All returns records
//   index of the field
//   val to find
func (g getter) All(
	index interface{},
	val interface{},
	offset uint32,
	limit uint32,
) ([]interface{}, error) {
	r, err := g.c.Select(g.s, index, offset, limit,
		tarantool.IterEq, []interface{}{val})
	if err != nil {
		return nil, err
	}
	records := make([]interface{}, 0)
	for _, t := range r.Data {
		if t, ok := t.([]interface{}); ok {
			m, err := g.model(t)
			if err != nil {
				return nil, err
			}
			records = append(records, m)
		} else {
			return nil, InvalidEntityError{}
		}
	}
	return records, nil
}

// FindBy record
// select record with a specific value
func (g getter) FindBy(index interface{}, val interface{}) (interface{}, error) {
	r, err := g.c.Select(g.s, index, 0, 1,
		tarantool.IterEq, []interface{}{val})
	if err != nil {
		return nil, err
	}
	if len(r.Data) == 0 {
		return nil, fmt.Errorf("Not found in the space: %s", g.s)
	}
	if t, ok := r.Data[0].([]interface{}); ok {
		return g.model(t)
	}
	return nil, InvalidEntityError{}
}

// Find a record by ID
func (g getter) Find(id string) (interface{}, error) {
	return g.FindBy(consts.PrimaryIndex, id)
}
