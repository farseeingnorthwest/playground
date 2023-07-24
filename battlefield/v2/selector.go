package battlefield

import "math/rand"

type Selector interface {
	Select(*Fighter, []*Fighter) *Fighter
}

type RandomSelector struct {
	Own bool
}

func (s *RandomSelector) Select(i *Fighter, fighters []*Fighter) *Fighter {
	var candidates []*Fighter
	for _, f := range fighters {
		if f.Functional() && f.Side == i.Side == s.Own {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	return candidates[int(rand.Float64()*float64(len(candidates)))]
}
