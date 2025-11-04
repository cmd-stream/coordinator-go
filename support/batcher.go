package support

import "fmt"

// TODO Refactor.

func NewBatcher(count int64) *Batcher {
	return &Batcher{count: count, chpt: -1}
}

func NewBatcherWithCheckpoint(count, chpt int64) *Batcher {
	return &Batcher{count: count, chpt: chpt}
}

type Batcher struct {
	count int64
	chpt  int64
	first *Interval
	last  *Interval
}

func (b *Batcher) Add(n int64) (chpt int64) {
	// TODO if n < b.chpt
	if b.first == nil {
		if n == b.chpt+1 && b.count == 1 {
			b.chpt = n
			return b.chpt
		}
		b.first = &Interval{Start: n, End: n}
		b.last = b.first
		return
	}

	if (b.increaseLeft(n, b.first) || b.increaseRight(n, b.first)) &&
		b.first.Len() >= b.count && b.first.Start == b.chpt+1 {
		b.chpt = b.first.End
		b.first = b.first.Next
		if b.first == b.last {
			b.last = nil
		}
		return b.chpt
	}

	curr := b.first.Next
	for curr != nil {
		if b.increaseLeft(n, curr) {
			return
		}
		if b.increaseRight(n, curr) {
			return
		}
		curr = curr.Next
	}

	if b.addFirst(n) {
		return
	}
	b.addLast(n)
	return
}

func (b *Batcher) String() (str string) {
	curr := b.first
	for curr != nil {
		str += curr.String()
		curr = curr.Next
	}
	return
}

func (b *Batcher) increaseLeft(n int64, intv *Interval) bool {
	if n == intv.Start-1 {
		if intv.Prev != nil && n == intv.Prev.End+1 {
			intv.Start = intv.Prev.Start

			if intv.Prev.Prev != nil {
				intv.Prev = intv.Prev.Prev
				intv.Prev.Prev.Next = intv
			} else {
				// intv.Prev == b.first
				intv.Prev = nil
				b.first = intv
			}
			return true
		}

		intv.Start = n
		return true
	}
	return false
}

func (b *Batcher) increaseRight(n int64, intv *Interval) bool {
	if n == intv.End+1 {
		if intv.Next != nil && n == intv.Next.Start-1 {
			intv.End = intv.Next.End

			if intv.Next.Next != nil {
				intv.Next = intv.Next.Next
				intv.Next.Next.Prev = intv
			} else {
				// intv.Next = b.last
				intv.Next = nil
				b.last = intv
			}
			return true
		}
		intv.End = n
		return true
	}
	return false
}

func (b *Batcher) addFirst(n int64) bool {
	if n < b.first.Start {
		newIntv := &Interval{
			Start: n,
			End:   n,
			Next:  b.first,
		}
		b.first.Prev = newIntv
		b.first = newIntv
		return true
	}
	return false
}

func (b *Batcher) addLast(n int64) bool {
	if n > b.last.End {
		newIntv := &Interval{
			Start: n,
			End:   n,
			Prev:  b.last,
		}
		b.last.Next = newIntv
		b.last = newIntv
		return true
	}
	return false
}

func NewInterval(start, end int64) *Interval {
	return &Interval{Start: start, End: end}
}

type Interval struct {
	Start int64
	End   int64

	Next *Interval
	Prev *Interval
}

func (i *Interval) Increase(n int64) bool {
	if n == i.Start-1 {
		i.Start = n
		return true
	}
	if n == i.End+1 {
		i.End = n
		return true
	}
	return false
}

func (i *Interval) Len() int64 {
	return i.End - i.Start + 1
}

func (i *Interval) String() string {
	return fmt.Sprintf("[%v, %v, %p, %p]", i.Start, i.End, i.Prev, i.Next)
}
