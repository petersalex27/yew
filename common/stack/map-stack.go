package stack

type MapStack[a comparable, b any] struct {
	*Stack[map[a]b]
}

func NewMap[a comparable, b any]() *MapStack[a, b] {
	return &MapStack[a, b]{
		Stack: New[map[a]b](),
	}
}

// maps a key-value pair to the top map
//
// if a mapping already exists for the key, it is overwritten
//
//   - SEE: Remap for a version that returns the old value if it exists
//   - SEE: RemapWith for a version that allows the new value to be computed from the old value
//   - SEE: MapNew for a version that only maps the key-value pair if the key does not already exist
func (s *MapStack[a, b]) Map(key a, value b) {
	if s.ctr == 0 {
		s.grow()
		s.Push(make(map[a]b))
	}

	s.data[s.ctr-1][key] = value
}

// maps a key-value pair to the top map, returning the old value (if it existed) and whether the old
// value existed
//
//   - SEE: Map for a version that does not return the old value
//   - SEE: RemapWith for a version that allows the new value to be computed from the old value
//   - SEE: MapNew for a version that only maps the key-value pair if the key does not already exist
func (s *MapStack[a, b]) Remap(key a, value b) (old b, found bool) {
	if s.ctr == 0 {
		return old, false
	}

	old, found = s.data[s.ctr-1][key]
	s.data[s.ctr-1][key] = value
	return old, found
}

// maps a key-value pair to the top map, allowing the new value to be computed from the old value
//
// if the key does not exist, the map is not modified and false is returned
//
//   - SEE: Map for a version that simply maps the key-value pair regardless of the current state
//   - SEE: Remap for a version that returns the old value if it exists
//   - SEE: MapNew for a version that only maps the key-value pair if the key does not already exist
func (s *MapStack[a, b]) RemapWith(key a, f func(b) b) (old b, found bool) {
	if s.ctr == 0 {
		return old, false
	}

	old, found = s.data[s.ctr-1][key]
	if found {
		s.data[s.ctr-1][key] = f(old)
	}
	return old, found
}

// un-maps a key-value pair from the top map
//
// if the key does not exist, the map is not modified and false is returned for `found`
func (s *MapStack[a, b]) Unmap(key a) (old b, found bool) {
	if s.ctr == 0 {
		return old, false
	}

	old, found = s.data[s.ctr-1][key]
	if found {
		delete(s.data[s.ctr-1], key)
	}
	return old, found
}

// maps a new key-value pair to the top map
//
// if the key already exists, the map is not modified and false is returned
//
//   - SEE: Map for a version that simply maps the key-value pair regardless of the current state
//   - SEE: Remap for a version that returns the old value if it exists, but maps the key-value pair
//     regardless of the current state
//   - SEE: RemapWith for a version that allows the new value to be computed from the old value
func (s *MapStack[a, b]) MapNew(key a, value b) (success bool) {
	if s.ctr == 0 {
		s.grow()
		s.Push(make(map[a]b))
	}

	if _, found := s.data[s.ctr-1][key]; found {
		return false
	}

	s.data[s.ctr-1][key] = value
	return true
}

// gets the value associated with the key in the top map and whether the key exists
func (s *MapStack[a, b]) Get(key a) (val b, found bool) {
	if s.ctr == 0 {
		return val, false
	}

	val, found = s.data[s.ctr-1][key]
	return val, found
}

// checks if the key exists in the top map
func (s *MapStack[a, b]) Exists(key a) (found bool) {
	_, found = s.Get(key)
	return found
}

// searches for the key in the maps starting from the top and going down
//
// if found, copies the key-value pair to the top map and returns true
//
// NOTE: this includes copying the key-value pair from the top map to the top map if found there.
//   - This is the first map that will be checked
//   - This is not the same as skipping the top map, since if it's only in the top map, this function
//     will return true unlike a version that skips the top map
//
// SEE:`CopyUpWith` for a version that allows the new value to be computed from the old value
func (s *MapStack[a, b]) CopyUp(key a) (success bool) {
	if s.ctr == 0 {
		return false
	}

	for i := s.ctr - 1; i >= 0; i-- {
		if val, found := s.data[i][key]; found {
			s.Map(key, val)
			return true
		}
	}
	return false
}

// Like `CopyUp`, but allows the new value to be computed from the old value
//
// SEE: CopyUp
func (s *MapStack[a, b]) CopyUpWith(key a, f func(b) b) (success bool) {
	if s.ctr == 0 {
		return false
	}

	for i := s.ctr - 1; i >= 0; i-- {
		if val, found := s.data[i][key]; found {
			s.Map(key, f(val))
			return true
		}
	}
	return false
}

func (s *MapStack[a, b]) filterHelper(predicate func(b) bool) (trueMap map[a]b, falseMap map[a]b) {
	if s.ctr == 0 {
		return make(map[a]b), make(map[a]b)
	}

	halfSize := len(s.data[s.ctr-1]) / 2
	trueMap = make(map[a]b, halfSize)
	falseMap = make(map[a]b, halfSize)

	for key, val := range s.data[s.ctr-1] {
		if predicate(val) {
			trueMap[key] = val
		} else {
			falseMap[key] = val
		}
	}

	return trueMap, falseMap
}

// removes all key-value pairs from the top map whose values satisfy the predicate
//
//   - SEE: Filter for a version that instead uses the predicate to decide which key-value pairs to
//     keep--still returning the removed key-value pairs
func (s *MapStack[a, b]) FilterOut(predicate func(b) bool) (out map[a]b) {
	out, keep := s.filterHelper(predicate)
	s.data[s.ctr-1] = keep
	return out
}

// removes all key-value pairs from the top map whose values do not satisfy the predicate
//
//   - SEE: FilterOut for a version that instead uses the predicate to decide which key-value pairs to
//     remove--still returning the removed key-value pairs
func (s *MapStack[a, b]) Filter(predicate func(b) bool) (out map[a]b) {
	keep, out := s.filterHelper(predicate)
	s.data[s.ctr-1] = keep
	return out
}

// removes the top map from the stack by folding its key-value pairs into a value and returning it
func (s *MapStack[a, b]) Fold(z b, f func(b, b) b) b {
	if s.ctr == 0 {
		return z
	}

	m, _ := s.Pop()
	for _, val := range m {
		z = f(z, val)
	}
	return z
}
