package getter

import (
	"fmt"
	"log"
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

// Creator is used for creating new model
type Creator interface {
	Create() interface{}
}

func (e InvalidEntityError) Error() string {
	return "Invalid entity "
}

// getter implements Getter
type getter struct {
	s string  // space
	m Creator // model
	c interfaces.IConnection
}

// Getter is a helper for records getting
type Getter interface {
	All(
		index interface{},
		val interface{},
		offset uint32,
		limit uint32,
	) ([]interface{}, error)
	FindBy(index interface{}, val interface{}) (interface{}, error)
	Find(id string) (interface{}, error)
	Count(index, val string) (int, error)
}

// NewGetter wraps a model with a getter methods
func NewGetter(
	space string,
	model Creator,
	conn interfaces.IConnection,
) Getter {
	return &getter{space, model, conn}
}

// model parses tuple into the model
func (g getter) model(t []interface{}) (interface{}, error) {
	var mdl = g.m.Create()
	var in = map[string]interface{}{}
	ref := reflect.TypeOf(mdl).Elem()

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

	if err := mapstructure.Decode(in, mdl); err != nil {
		return nil, err
	}

	return mdl, nil
}

func (g getter) Count(index, val string) (int, error) {
	funcName := fmt.Sprint("box.space.", g.s, ".index.", index, ":count")
	resp, err := g.c.Call(
		funcName,
		[]interface{}{
			val,
			map[string]string{
				"iterator": "EQ",
			},
		})
	log.Print(resp)
	return 0, err
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
		tarantool.IterReq, []interface{}{val})
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
