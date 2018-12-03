package arrlist

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var debugMode = false

func TestEmptyArrList(t *testing.T) {
	l := New(0)
	require.Equal(t, 0, l.Len())

	c := 0
	l.ForEach(func(v interface{}) bool {
		c++
		return false
	})
	require.Equal(t, 0, c)
}

func TestArrListPush(t *testing.T) {
	l := New(10)
	for i := 0; i < 20; i++ {
		l.Push(i)
		dump(fmt.Sprintf("push %d", i), l)
	}
	require.Equal(t, 20, l.Len())

	c := 0
	l.ForEach(func(v interface{}) bool {
		require.Equal(t, c, v)
		c++
		return false
	})
	require.Equal(t, 20, c)
}

func TestArrListRemove(t *testing.T) {
	l := New(5)

	ids := make([]*interface{}, 12)
	for i := range ids {
		ids[i] = l.Push(i)
		dump(fmt.Sprintf("push %d", i), l)
	}

	idToInt := func(id *interface{}) int {
		v := *id
		return v.(int)
	}

	removed := []int{} // 已经删除了哪些值
	check := func() {
		msg := fmt.Sprintf("check(removed: %+v)", removed)
		require.Equal(t, len(ids), l.Len(), msg)

		expected := make([]int, len(ids))
		for i, id := range ids {
			expected[i] = idToInt(id)
		}

		actual := []int{}
		l.ForEach(func(v interface{}) bool {
			actual = append(actual, v.(int))
			return false
		})
		require.EqualValues(t, expected, actual, msg)
	}

	check()

	for i := 0; i < len(ids); i++ {
		index := rand.Intn(len(ids))

		id := ids[index]
		v := idToInt(id)

		// remove from ids
		ids = append(ids[:index], ids[index+1:]...)
		removed = append(removed, v)

		l.Remove(id)
		dump(fmt.Sprintf("remove %v", v), l)
		check()
	}
}

func TestArrListForEach(t *testing.T) {
	l := New(0)
	for i := 0; i < 10; i++ {
		l.Push(i)
	}

	values := make([]int, 0)
	l.ForEach(func(v interface{}) bool {
		values = append(values, v.(int))
		return len(values) >= 5
	})
	require.Len(t, values, 5)
	for i, v := range values {
		require.Equal(t, i, v)
	}
}

func dump(header string, l ArrList) {
	if !debugMode {
		return
	}
	os.Stdout.WriteString(header)
	os.Stdout.WriteString("\n")
	if al, ok := l.(*arrList); ok {
		al.dumpPoints(os.Stdout)
	}
	os.Stdout.WriteString("\n")
}

func (l *arrList) dumpPoints(w io.Writer) {
	fmt.Fprintf(w, "freeHead = %p\n", l.freeHead)
	for i, page := range l.pagePool {
		for j, p := range page {
			fmt.Fprintf(w, "page[%d][%d](%T) = %v(p: %p)\n", i, j, p, p, &page[j])
		}
	}
}
