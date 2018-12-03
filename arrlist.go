package arrlist

// DefaultPageSize 默认每页的大小
const DefaultPageSize = 1024

// ArrList 以数组的形式实现单链表。
type ArrList interface {
	// Push 将 v 添加到链表中，返回的 id 用于从链表中删除该元素。
	Push(v interface{}) (id *interface{})
	// Remove 从链表中删除 id 指向的元素。
	// Remove 之后 id 就不能再使用了。
	Remove(id *interface{})
	// ForEach 遍历链表，对每个元素都调用 f，如果 stop 为 true 则停止遍历
	ForEach(f func(v interface{}) (stop bool))
	// Len 返回当前链表中已存在多少元素
	Len() int
}

var _ ArrList = (*arrList)(nil)

type arrList struct {
	// pagePool 中每个页都有 pageSize 大小，在初始时一次性分配。
	pagePool [][]interface{}
	// freeHead 指向下一个可用的地址
	freeHead *interface{}
	// pageSize 指定每页的长度
	pageSize int
	// length 表示当前链表中有多少元素
	length int
}

// New 返回 ArrList 的一个默认实现。
func New(pageSize int) ArrList {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return &arrList{
		pageSize: pageSize,
	}
}

func (l *arrList) Push(v interface{}) (id *interface{}) {
	if l.freeHead == nil {
		l.appendPage()
	}
	id = l.freeHead
	// 指向下一个可用地址，(*freeHead) 中保存的是下一个可用的元素地址
	l.freeHead = (*l.freeHead).(*interface{})
	*id = v
	l.length++
	return
}

func (l *arrList) Remove(id *interface{}) {
	// 在 id 中保存下一个可用的元素地址
	*id = l.freeHead
	// 使 id 成为下一个可用的元素地址
	l.freeHead = id
	l.length--
}

func (l *arrList) ForEach(f func(v interface{}) (stop bool)) {
	for _, page := range l.pagePool {
		for _, v := range page {
			// 当类型是 *interface{} 时表明该位置还没使用
			if _, ok := v.(*interface{}); !ok {
				if stop := f(v); stop {
					return
				}
			}
		}
	}
}

func (l *arrList) Len() int { return l.length }

func (l *arrList) appendPage() {
	page := make([]interface{}, l.pageSize)
	// 前一个位置保存后一个位置的地址，最后一个位置保存 nil。
	for i := 1; i < len(page); i++ {
		page[i-1] = &page[i]
	}
	// 这句不能少，赋值前的类型是 nil，值也是 nil，赋值后的类型是 *interface{}，
	// 后续会有类型强转，所以这里不能省略。
	page[len(page)-1] = l.freeHead
	l.pagePool = append(l.pagePool, page)
	l.freeHead = &page[0]
}
