// Package mutation has methods for entities changing
package mutation

import (
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/hedhyw/go-tarantool-structure/internal/consts"
	"github.com/hedhyw/go-tarantool-structure/internal/errs"
	db "github.com/hedhyw/go-tarantool-structure/internal/interfaces"
)

// NewMutation wraps a model with a mutation methods
func NewMutation(
	space string, model interface{},
	conn db.IConnection,
) Mutation {
	return &mutation{space, model, conn}
}

// mutation implements Mutation
type mutation struct {
	s string      // space
	m interface{} // model
	c db.IConnection
}

// Mutation has methods to operate with entity in the database
type Mutation interface {
	Append() (id string, err error)
	Update() error
	Delete() error
}

// id of the record
func (m mutation) id() interface{} {
	ref := reflect.TypeOf(m.m).Elem()
	for i, count := 0, ref.NumField(); i < count; i++ {
		f := ref.Field(i)
		val, ok := f.Tag.Lookup(consts.TagName)
		if ok && val == "0" {
			return reflect.ValueOf(m.m).Elem().Field(i).Interface()
		}
	}
	return nil
}

func tupleElements(
	refV reflect.Value,
	refT reflect.Type,
	t *map[int]interface{},
) (lastIndex int, err error) {
	count := refT.NumField()
	for i := 0; i < count; i++ {
		f := refT.Field(i)
		var index int
		v := refV.Field(i)
		if v.Kind() == reflect.Struct {
			index, err = tupleElements(v, v.Type(), t)
		} else {
			val, ok := f.Tag.Lookup(consts.TagName)
			if !ok {
				continue
			}
			index, err = strconv.Atoi(val)
			if err != nil {
				return 0, errs.NewInvalidTagValueError(err.Error())
			}
			if refV.Field(i).CanInterface() {
				(*t)[index] = refV.Field(i).Interface()
			}
		}
		if index > lastIndex {
			lastIndex = index
		}
	}
	return
}

// tuple is a raw result
func (m mutation) tuple() (tuple []interface{}, err error) {
	t := make(map[int]interface{})
	refV := reflect.ValueOf(m.m).Elem()
	refT := reflect.TypeOf(m.m).Elem()
	lastIndex, err := tupleElements(refV, refT, &t)
	if err != nil {
		return
	}
	tuple = make([]interface{}, lastIndex+1)
	for ind, val := range t {
		tuple[ind] = val
	}
	return tuple, nil
}

// Append adds an entity to the database
func (m mutation) Append() (id string, err error) {
	id = uuid.New().String()
	t, err := m.tuple()
	if err != nil {
		return "", err
	}
	t[0] = id
	_, err = m.c.Insert(m.s, t)
	return id, err
}

// Update entity with new values
func (m mutation) Update() error {
	t, err := m.tuple()
	if err != nil {
		return err
	}
	_, err = m.c.Replace(m.s, t)
	return err
}

// Delete entity
func (m mutation) Delete() error {
	_, err := m.c.Delete(m.s, consts.PrimaryIndex, []interface{}{m.id()})
	return err
}
