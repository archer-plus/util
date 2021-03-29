package redix

// Iterator 迭代器
type Iterator struct {
	data  []interface{}
	index int
}

// NewIterator 创建新的迭代器
func NewIterator(data []interface{}) *Iterator {
	return &Iterator{data: data}
}

// HasNext 是否有下一个值
func (i *Iterator) HasNext() bool {
	if i.data == nil || len(i.data) == 0 {
		return false
	}
	return i.index < len(i.data)
}

// Next 返回下一个值
func (i *Iterator) Next() (res interface{}) {
	res = i.data[i.index]
	i.index++
	return
}
