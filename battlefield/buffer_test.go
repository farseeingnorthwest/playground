package battlefield

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferList_Len(t *testing.T) {
	for _, tt := range []struct {
		l bufferList
		n int
	}{
		{bufferList{}, 0},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
			}},
			1,
		},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
				next: &bufferNode{
					Buffer: healthCriticalBaseline{},
				},
			}},
			2,
		},
	} {
		t.Run(fmt.Sprintf("%v", tt.l), func(t *testing.T) {
			assert.Equal(t, tt.n, tt.l.Len())
		})
	}
}

func TestBufferList_Append(t *testing.T) {
	for i, tt := range []struct {
		l bufferList
		b Buffer
		n int
	}{
		{bufferList{}, healthCriticalBaseline{}, 1},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
			}},
			healthCriticalBaseline{},
			2,
		},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
				next: &bufferNode{
					Buffer: healthCriticalBaseline{},
				},
			}},
			healthCriticalBaseline{},
			3,
		},
	} {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			tt.l.Append(tt.b)
			p := (*bufferNode)(&tt.l)
			for i := 0; i < tt.n; i++ {
				p = p.next
			}
			assert.Equal(t, tt.b, p.Buffer)
			assert.Nil(t, p.next)
		})
	}
}

func TestBufferList_Drain(t *testing.T) {
	buffers := bufferList{
		next: &bufferNode{
			Buffer: healthCriticalBaseline{},
			next: &bufferNode{
				Buffer: &volatileBuffer{},
				next: &bufferNode{
					Buffer: &volatileBuffer{r: 2},
					next: &bufferNode{
						Buffer: &volatileBuffer{r: 3},
						next: &bufferNode{
							Buffer: &volatileBuffer{r: 5},
						},
					},
				},
			},
		}}

	for _, n := range []int{4, 4, 3, 2, 2, 1, 1} {
		buffers.Drain()
		assert.Equal(t, n, buffers.Len())
	}
}

func TestBufferList_Buff(t *testing.T) {
	for i, tt := range []struct {
		l bufferList
		v float64
	}{
		{bufferList{}, 1},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
			}},
			1.5,
		},
		{bufferList{
			next: &bufferNode{
				Buffer: healthCriticalBaseline{},
				next: &bufferNode{
					Buffer: healthCriticalBaseline{},
				},
			}},
			2.25,
		},
	} {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			assert.Equal(t, tt.v, tt.l.Buff(1))
		})
	}
}

type volatileBuffer struct {
	r int
	v float64
}

func (b *volatileBuffer) Drain() int {
	r := b.r
	b.r--

	return r
}

func (b *volatileBuffer) Buff(float64) float64 {
	return b.v
}
