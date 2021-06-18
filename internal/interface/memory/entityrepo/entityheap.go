package entityrepo

type MusicalEntityItem struct {
	id        string
	checkedAt uint64
}

type MusicalEntityHeap []MusicalEntityItem

func (h MusicalEntityHeap) Len() int           { return len(h) }
func (h MusicalEntityHeap) Less(i, j int) bool { return h[i].checkedAt < h[j].checkedAt }
func (h MusicalEntityHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MusicalEntityHeap) Push(x interface{}) {
	item := x.(MusicalEntityItem)
	*h = append(*h, item)
}

func (h *MusicalEntityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
