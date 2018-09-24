package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testModel struct {
	Number int64  `tnt:"1"`
	ID     string `tnt:"0"`
}

var testTuple = []interface{}{"XQE0ZYBSYP", int64(10)}

const testSpace = "TestMSpace"

func model() testModel {
	return testModel{
		ID:     testTuple[0].(string),
		Number: testTuple[1].(int64),
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
