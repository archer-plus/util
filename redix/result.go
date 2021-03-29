package redix

// Result 查询结果
type Result struct {
	Value interface{}
	Error error
}

// NewResult 新建查询结果
func NewResult(value interface{}, error error) *Result {
	return &Result{Value: value, Error: error}
}

// Result 返回查询结果的值
func (c *Result) Result(v ...interface{}) interface{} {
	if c.Error != nil {
		if len(v) > 0 {
			return v[0]
		}
		panic(c.Error)
	}
	return c.Value
}

// MResult 列表查询结果
type MResult struct {
	Value []interface{}
	Error error
}

// NewMResult 新建列表查询结果
func NewMResult(value []interface{}, error error) *MResult {
	return &MResult{Value: value, Error: error}
}

// Result 返回列表查询结果的值
func (m *MResult) Result(v []interface{}) []interface{} {
	if m.Error != nil {
		return v
	}
	return m.Value
}

// Iterator 列表迭代器
func (m *MResult) Iterator() *Iterator {
	return NewIterator(m.Value)
}
