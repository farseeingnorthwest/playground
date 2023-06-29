package battlefield

type Buffer interface {
	Drain() int
	Buff(float64) float64
}

type healthCriticalBaseline struct{}

func (healthCriticalBaseline) Drain() int {
	return 1
}

func (healthCriticalBaseline) Buff(value float64) float64 {
	return value * 1.5
}

type bufferNode struct {
	Buffer
	next *bufferNode
}

type bufferList bufferNode

func (l *bufferList) Len() int {
	p := (*bufferNode)(l)
	n := 0
	for p.next != nil {
		n++
		p = p.next
	}

	return n
}

func (l *bufferList) Append(buffer Buffer) {
	p := (*bufferNode)(l)
	for p.next != nil {
		p = p.next
	}
	p.next = &bufferNode{buffer, nil}
}

func (l *bufferList) Drain() {
	p := (*bufferNode)(l)
	for p.next != nil {
		if p.next.Drain() <= 0 {
			p.next = p.next.next
		} else {
			p = p.next
		}
	}
}

func (l *bufferList) Buff(value float64) float64 {
	p := (*bufferNode)(l)
	for p.next != nil {
		value = p.next.Buff(value)
		p = p.next
	}

	return value
}
