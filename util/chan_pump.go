package util

type Pump struct {
	Out   chan any
	In    map[int]<-chan any
	index int
}

func NewPump() *Pump {
	return &Pump{
		Out:   make(chan any, 1),
		In:    make(map[int]<-chan any, 0),
		index: -1,
	}
}

func (p *Pump) AddIn(in <-chan any) {
	p.index++
	p.In[p.index] = in
}

func (p *Pump) RunForward() {
	msgSent := make(chan struct{})
	for idx, in := range p.In {
		go func(in <-chan any, idx int) {
			select {
			case <-msgSent:
				return
			case <-in:
				close(msgSent)
				p.Out <- idx
				return
			}
		}(in, idx)
	}
}
