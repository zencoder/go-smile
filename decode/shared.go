package decode

import "fmt"

type SharedState struct {
	sharedKeyNames []interface{}
	sharedValues   []interface{}
}

func (s *SharedState) AddSharedValue(sharedValue interface{}) {
	if len(s.sharedValues) == 1024 {
		s.sharedValues = []interface{}{}
	}
	s.sharedValues = append(s.sharedValues, sharedValue)
}

func (s *SharedState) GetSharedValue(index int) (interface{}, error) {
	if index >= len(s.sharedValues) {
		return nil, fmt.Errorf("shared value %d requested but only %d values available", index, len(s.sharedValues))
	}
	return s.sharedValues[index], nil
}

func (s *SharedState) AddSharedKey(sharedKeyName interface{}) {
	if len(s.sharedKeyNames) == 1024 {
		s.sharedKeyNames = []interface{}{}
	}
	s.sharedKeyNames = append(s.sharedKeyNames, sharedKeyName)
}

func (s *SharedState) GetSharedKey(index int) (interface{}, error) {
	if index >= len(s.sharedKeyNames) {
		return nil, fmt.Errorf("shared key %d requested but only %d keys available", index, len(s.sharedKeyNames))
	}
	return s.sharedKeyNames[index], nil
}
