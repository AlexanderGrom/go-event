package event

import (
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOn(t *testing.T) {
	testCases := []struct {
		name      string
		fn        interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "test.event.on.1",
			fn:        func() error { return nil },
			assertion: assert.NoError,
		}, {
			name:      "test.event.on.2",
			fn:        func(i int) error { return nil },
			assertion: assert.NoError,
		}, {
			name:      "test.event.on.2",
			fn:        func(i int) error { return nil },
			assertion: assert.NoError,
		}, {
			name:      "test.event.on.2",
			fn:        func(i string) error { return nil },
			assertion: assert.Error,
		}, {
			name:      "test.event.on.2",
			fn:        func(i int, j string) error { return nil },
			assertion: assert.Error,
		}, {
			name:      "test.event.on.3",
			fn:        func() int { return 0 },
			assertion: assert.Error,
		}, {
			name:      "test.event.on.4",
			fn:        func() (int, error) { return 0, nil },
			assertion: assert.Error,
		}, {
			name:      "test.event.on.5",
			fn:        func() {},
			assertion: assert.Error,
		}, {
			name:      "test.event.on.6",
			fn:        nil,
			assertion: assert.Error,
		}, {
			name:      "test.event.on.7",
			fn:        "func",
			assertion: assert.Error,
		}, {
			name:      "test.event.on.8",
			fn:        func(e Event) error { return nil },
			assertion: assert.NoError,
		},
	}

	e := New()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := e.On(tt.name, tt.fn)
			tt.assertion(t, err)
		})
	}
}

func TestGo(t *testing.T) {
	type str struct {
		count int
	}
	counter := 0
	testCases := []struct {
		name      string
		fn        interface{}
		params    []interface{}
		count     int
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "test.event.go.1",
			fn:        func() error { counter++; return nil },
			params:    []interface{}{},
			count:     1,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.1",
			fn:        func() error { counter++; return nil },
			params:    []interface{}{},
			count:     3,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.2",
			fn:        func() error { counter = 0; return nil },
			params:    []interface{}{},
			count:     0,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.3",
			fn:        func() error { return errors.New("some error") },
			params:    []interface{}{},
			count:     0,
			assertion: assert.Error,
		}, {
			name:      "test.event.go.4",
			fn:        func(i int) error { counter = i; return nil },
			params:    []interface{}{1},
			count:     1,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.5",
			fn:        func(i int, j int) error { counter = i + j; return nil },
			params:    []interface{}{1, 2},
			count:     3,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.6",
			fn:        func(i int, j ...int) error { counter = len(j); return nil },
			params:    []interface{}{1, 2, 3},
			count:     2,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.7",
			fn:        func(i string, j ...int) error { counter = len(j); return nil },
			params:    []interface{}{"A", 1, 2},
			count:     2,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.8",
			fn:        func(i str, j *str, k ...int) error { counter = len(k) + j.count; return nil },
			params:    []interface{}{str{1}, &str{2}, 3, 4},
			count:     4,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.9",
			fn:        func(i ...int) error { counter = len(i); return nil },
			params:    []interface{}{},
			count:     0,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.10",
			fn:        func(e Event) error { return nil },
			params:    []interface{}{nil},
			count:     0,
			assertion: assert.NoError,
		}, {
			name:      "test.event.go.11",
			fn:        func(e Event) error { return nil },
			params:    []interface{}{},
			count:     0,
			assertion: assert.Error,
		}, {
			name:      "test.event.go.12",
			fn:        func(i int, j ...int) error { return nil },
			params:    []interface{}{},
			count:     0,
			assertion: assert.Error,
		}, {
			name:      "test.event.go.13",
			fn:        func(i int) error { return nil },
			params:    []interface{}{},
			count:     0,
			assertion: assert.Error,
		},
		// Добавить проверки на отсутствующие параметры 3 штуки!!!
	}

	e := New()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := e.On(tt.name, tt.fn)
			assert.NoError(t, err)

			err = e.Go(tt.name, tt.params...)
			tt.assertion(t, err)

			assert.Equal(t, tt.count, counter)
		})
	}
}

func TestHas(t *testing.T) {
	e := New()
	e.On("test.event.has.1", func() error { return nil })
	e.On("test.event.has.2", func() error { return nil })
	assert.True(t, e.Has("test.event.has.1"))
	assert.False(t, e.Has("test.event.has.3"))
}

func TestList(t *testing.T) {
	e := New()
	e.On("test.event.ls.1", func() error { return nil })
	e.On("test.event.ls.1", func() error { return nil })
	e.On("test.event.ls.2", func() error { return nil })
	e.On("test.event.ls.3", func() error { return nil })
	list := e.List()
	assert.Equal(t, 3, len(list))
	sort.Strings(list)
	assert.Equal(t, []string{"test.event.ls.1", "test.event.ls.2", "test.event.ls.3"}, list)
}

func TestRemove(t *testing.T) {
	e := New()
	e.On("test.event.has.1", func() error { return nil })
	e.On("test.event.has.2", func() error { return nil })
	e.On("test.event.has.3", func() error { return nil })
	assert.Equal(t, 3, len(e.List()))
	e.Remove("test.event.has.2")
	assert.Equal(t, 2, len(e.List()))
	assert.False(t, e.Has("test.event.has.2"))
	e.Remove("test.event.has.2")
	assert.Equal(t, 2, len(e.List()))
	e.Remove()
	assert.Equal(t, 0, len(e.List()))
}
