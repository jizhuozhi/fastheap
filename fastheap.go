package fastheap

const minCap = 16

type LessFunc func(a, b interface{}) bool

// Heap a faster implementation than golang's heap
type Heap struct {
	elements []interface{}
	size     int
	less     LessFunc
}

func New(less LessFunc) *Heap {
	return NewWithCap(minCap, less)
}

func NewWithCap(cap int, less LessFunc) *Heap {
	if cap < minCap {
		cap = minCap
	}
	return &Heap{
		elements: make([]interface{}, cap),
		less:     less,
	}
}

func (h *Heap) Fix(i int) {
	if !h.fixDown(i, h.size) {
		h.fixUp(i)
	}
}

func (h *Heap) Push(element interface{}) {
	h.ensureIncrement()
	n := h.size
	h.size++
	h.elements[n] = element
	h.fixUp(n)
}

func (h *Heap) Peek() interface{} {
	return h.elements[0]
}

func (h *Heap) Pop() interface{} {
	value := h.elements[0]
	h.size--
	n := h.size
	h.elements[0] = h.elements[n]
	h.elements[n] = nil // For gc
	h.fixDown(0, n)
	h.ensureDecrement() // Avoid only growing but not decreasing
	return value
}

func (h *Heap) Size() int {
	return h.size
}

func (h *Heap) Empty() bool {
	return h.size == 0
}

func (h *Heap) ensureIncrement() {
	if h.size+1 > cap(h.elements) {
		oldElements := h.elements
		h.elements = make([]interface{}, 2*cap(h.elements))
		copy(h.elements, oldElements)
	}
}

func (h *Heap) ensureDecrement() {
	if minCap < cap(h.elements) && h.size*2 < cap(h.elements) {
		newCap := cap(h.elements) / 2
		oldElements := h.elements
		h.elements = make([]interface{}, newCap)
		copy(h.elements, oldElements)
	}
}

func (h *Heap) fixUp(i int) bool {
	oldI := i
	var parent int
	var element interface{}
	for element = h.elements[i]; i > 0; i = parent {
		parent = (i - 1) / 2
		if h.less(element, h.elements[parent]) {
			h.elements[i] = h.elements[parent]
		} else {
			break
		}
	}
	if oldI == i {
		return false
	}
	h.elements[i] = element
	return true
}

func (h *Heap) fixDown(i, n int) bool {
	oldI := i
	var child int
	var element interface{}
	for element = h.elements[i]; ; i = child {
		child = i*2 + 1
		if child < n {
			if child+1 < n && h.less(h.elements[child+1], h.elements[child]) {
				child++
			}
			if h.less(h.elements[child], element) {
				h.elements[i] = h.elements[child]
			} else {
				break
			}
		} else {
			break
		}
	}
	if oldI == i {
		return false
	}
	h.elements[i] = element
	return true
}
