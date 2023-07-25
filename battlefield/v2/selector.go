package battlefield

import "math/rand"

type Selector interface {
	Select(Warrior, []Warrior) Warrior
}

type RandomSelector struct {
	Own bool
}

func (s *RandomSelector) Select(i Warrior, fighters []Warrior) Warrior {
	var candidates []Warrior
	for _, f := range fighters {
		if f.Health() > 0 && f.(*Fighter).Side == i.(*Fighter).Side == s.Own {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	return candidates[int(rand.Float64()*float64(len(candidates)))]
}
