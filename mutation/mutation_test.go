package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testModel struct {
	service1 int32
	service2 int64
	Number   int64  `tnt:"1"`
	ID       string `tnt:"0"`
	Deep     deepTest
}

type deepTest struct {
	Field1 int  `tnt:"2"`
	Field2 bool `tnt:"3"`
}

var testTuple = []interface{}{"XQE0ZYBSYP", int64(10), int(8), true}

const testSpace = "TestMSpace"

func model() testModel {
	return testModel{
		ID:     testTuple[0].(string),
		Number: testTuple[1].(int64),
		Deep: deepTest{
			Field1: testTuple[2].(int),
			Field2: testTuple[3].(bool),
		},
	}
}

func TestTuple(t *testing.T) {
	m := model()
	mutRaw := NewMutation(testSpace, &m, nil)
	mut, ok := mutRaw.(*mutation)
	assert.True(t, ok)
	tup, err := mut.tuple()
	assert.Nil(t, err)
	assert.Equal(t, testTuple, tup)
}

func TestID(t *testing.T) {
	m := model()
	mutRaw := NewMutation(testSpace, &m, nil)
	mut, ok := mutRaw.(*mutation)
	assert.True(t, ok)
	assert.Equal(t, testTuple[0], mut.id())
}
