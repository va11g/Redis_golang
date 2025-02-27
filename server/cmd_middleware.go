package server

import "reflect"

type filterFunc = func([][]byte) ([][]byte, error)

type middleware struct {
	filters []filterFunc
}

func NewMiddleWare() *middleware {
	return &middleware{filters: make([]filterFunc, 0)}
}

func (m *middleware) Add(f filterFunc) {
	m.filters = append(m.filters, f)
}

func (m *middleware) Delete(f filterFunc) bool {
	for k, v := range m.filters {
		if reflect.ValueOf(f).Pointer() == reflect.ValueOf(v).Pointer() {
			m.filters = append(m.filters[:k], m.filters[k+1:]...)
			return true
		}
	}
	return false
}

func (m *middleware) Filter(cmd [][]byte) ([][]byte, error) {
	var err error
	for _, v := range m.filters {
		cmd, err = v(cmd)
		if err != nil {
			return nil, err
		}
	}
	return cmd, nil
}
