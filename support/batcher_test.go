package support

import (
	"testing"

	assert "github.com/ymz-ncnk/assert/panic"
)

func TestSeqs(t *testing.T) {
	t.Run("count == 1", func(t *testing.T) {
		seqs := NewBatcher(1)

		var chpt int64
		chpt = seqs.Add(1)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(2)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 3)
	})

	t.Run("count == n, sequential", func(t *testing.T) {
		seqs := NewBatcher(3)

		var chpt int64
		chpt = seqs.Add(1)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(2)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 3)
		chpt = seqs.Add(4)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(5)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(6)
		assert.Equal(chpt, 6)
	})

	t.Run("count == n, not sequential", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(2)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(1)
		assert.Equal(chpt, 2)
	})

	t.Run("count == n, two intervals, increase left", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(4)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(2)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(1)
		assert.Equal(chpt, 4)
	})

	t.Run("count == n, two intervals, increase right", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(3)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(4)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(1)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(2)
		assert.Equal(chpt, 4)
	})

	t.Run("count == n, three intervals, increase left", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(6)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(4)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(2)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(5)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(1)
		assert.Equal(chpt, 6)
	})

	t.Run("count == n, three intervals, increase right", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(5)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(1)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(6)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(4)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(4)
		assert.Equal(chpt, 6)
	})

	t.Run("add last", func(t *testing.T) {
		seqs := NewBatcher(2)

		var chpt int64
		chpt = seqs.Add(1)
		assert.Equal(chpt, 0)
		chpt = seqs.Add(3)
		assert.Equal(chpt, 0)

		assert.Equal(seqs.last.End, 3)
		assert.Equal(seqs.last.Next, nil)
	})
}
