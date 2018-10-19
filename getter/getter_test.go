package getter

/*

		TODO:
- count() test

*/

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hedhyw/go-tarantool-structure/internal/interfaces"
	"github.com/hedhyw/go-tarantool-structure/internal/testdata/mock_interfaces"
	"github.com/stretchr/testify/assert"
	tarantool "github.com/tarantool/go-tarantool"
)

type testModel struct {
	Number int64  `tnt:"1"`
	ID     string `tnt:"0"`
}

func (testModel) Create() interface{} {
	return &testModel{}
}

var testTuple = []interface{}{"CYT2CV0LYQ", int64(15)}

const testSpace = "TestGSpace"

func model() testModel {
	return testModel{
		ID:     testTuple[0].(string),
		Number: testTuple[1].(int64),
	}
}

func TestModel(t *testing.T) {
	m := model()
	gRaw := NewGetter(testSpace, &m, nil)
	g, ok := gRaw.(*getter)
	assert.True(t, ok)
	parsed, err := g.model(testTuple)
	assert.Nil(t, err)
	if parsed, ok := parsed.(*testModel); ok {
		assert.Equal(t, testTuple[0], parsed.ID)
		assert.Equal(t, testTuple[1], parsed.Number)
	} else {
		assert.Fail(t, "Parsed model is not testModel")
	}
}

func TestAll(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	var c = mock_interfaces.NewMockIConnection(ctl)
	in := []interface{}{testTuple, testTuple}
	res := tarantool.Response{
		Data: in,
	}

	c.EXPECT().Select(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&res, nil).Times(1)

	m := testModel{}
	g := NewGetter(testSpace, &m, interfaces.IConnection(c))
	records, err := g.All(1, in[1], 0, uint32(len(in)))
	assert.Nil(t, err)
	assert.Equal(t, len(in), len(records))
	for _, r := range records {
		v := r.(*testModel)
		assert.Equal(t, testTuple[0], v.ID)
		assert.Equal(t, testTuple[1], v.Number)
	}
}

func TestFind(t *testing.T) {
	var testID = testTuple[0].(string)

	ctl := gomock.NewController(t)
	defer ctl.Finish()
	var c = mock_interfaces.NewMockIConnection(ctl)
	res := tarantool.Response{
		Data: []interface{}{testTuple},
	}

	c.EXPECT().Select(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		[]interface{}{testID},
	).Return(&res, nil).Times(1)

	m := testModel{}
	g := NewGetter(testSpace, &m, interfaces.IConnection(c))
	r, err := g.Find(testID)
	tm := r.(*testModel)

	assert.Nil(t, err)
	assert.Equal(t, testID, tm.ID)
	assert.Equal(t, testTuple[1], tm.Number)
}
