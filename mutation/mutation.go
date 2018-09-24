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

// tuple is a raw result
func (m mutation) tuple() ([]interface{}, error) {
	refV := reflect.ValueOf(m.m).Elem()
	refT := reflect.TypeOf(m.m).Elem()
	count := refT.NumField()
	t := make([]interface{}, count)
	lastIndex := 0
	for i := 0; i < count; i++ {
		f := refT.Field(i)
		val, ok := f.Tag.Lookup(consts.TagName)
		if !ok {
			continue
		}
		index, err := strconv.Atoi(val)
		if err != nil || index >= count {
			return nil, errs.NewInvalidTagValueError(f.Name)
		}
		if index > lastIndex {
			lastIndex = index
		}
		t[index] = refV.Field(i).Interface()
	}
	return t[:lastIndex+1], nil
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
