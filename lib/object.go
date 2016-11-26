package jd

import (
	"fmt"
	"reflect"
	"sort"
)

type jsonObject map[string]JsonNode

var _ JsonNode = jsonObject(nil)

func (s jsonObject) Json() string {
	return renderJson(s)
}

func (s1 jsonObject) Equals(n JsonNode) bool {
	s2, ok := n.(jsonObject)
	if !ok {
		return false
	}
	if len(s1) != len(s2) {
		return false
	}
	return reflect.DeepEqual(s1, s2)
}

func (s1 jsonObject) Diff(n JsonNode) Diff {
	return s1.diff(n, Path{})
}

func (s1 jsonObject) diff(n JsonNode, path Path) Diff {
	d := make(Diff, 0)
	s2, ok := n.(jsonObject)
	if !ok {
		// Different types
		e := DiffElement{
			Path:     path.clone(),
			OldValue: s1,
			NewValue: n,
		}
		return append(d, e)
	}
	s1Keys := make([]string, 0, len(s1))
	for k := range s1 {
		s1Keys = append(s1Keys, k)
	}
	sort.Strings(s1Keys)
	s2Keys := make([]string, 0, len(s2))
	for k := range s2 {
		s2Keys = append(s2Keys, k)
	}
	sort.Strings(s2Keys)
	for _, k1 := range s1Keys {
		v1 := s1[k1]
		if v2, ok := s2[k1]; ok {
			// Both keys are present
			subDiff := v1.diff(v2, append(path.clone(), k1))
			d = append(d, subDiff...)
		} else {
			// S2 missing key
			e := DiffElement{
				Path:     append(path.clone(), k1),
				OldValue: v1,
				NewValue: voidNode{},
			}
			d = append(d, e)
		}
	}
	for _, k2 := range s2Keys {
		v2 := s2[k2]
		if _, ok := s1[k2]; !ok {
			// S1 missing key
			e := DiffElement{
				Path:     append(path.clone(), k2),
				OldValue: voidNode{},
				NewValue: v2,
			}
			d = append(d, e)
		}
	}
	return d
}

func (s jsonObject) Patch(d Diff) (JsonNode, error) {
	return patchAll(s, d)
}

func (o jsonObject) patch(pathBehind, pathAhead Path, oldValue, newValue JsonNode) (JsonNode, error) {
	// Base case
	if len(pathAhead) == 0 {
		if !o.Equals(oldValue) {
			return nil, fmt.Errorf(
				"Found %v at %v. Expected %v.",
				o.Json(), pathBehind, oldValue.Json())
		}
	}
	// Recursive case
	pe, ok := pathAhead[0].(string)
	if !ok {
		return nil, fmt.Errorf(
			"Found %v at %v. Expected JSON object.",
			o.Json(), pathBehind)
	}
	nextNode, ok := o[pe]
	if !ok {
		nextNode = voidNode{}
	}
	patchedNode, err := nextNode.patch(append(pathBehind, pe), pathAhead[1:], oldValue, newValue)
	if err != nil {
		return nil, err
	}
	if isVoid(patchedNode) {
		// Delete a pair
		delete(o, pe)
	} else {
		// Add or replace a pair
		o[pe] = patchedNode
	}
	return o, nil
}
