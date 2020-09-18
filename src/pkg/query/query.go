package query

// Query
type Query struct {
	PageNum  int
	PageSize int
	Cond     string
	Values   []interface{}
	Total    int64
}

// NewQuery
func NewQuery(pageNum int, pageSize int) *Query {
	if pageNum == 0 {
		pageNum = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}
	q := &Query{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    0,
	}

	q.PageNum = (q.PageNum - 1) * q.PageSize

	return q
}

func (q *Query) ValidCond(cond map[string]interface{}) (ret *Query, err error) {
	if q.Cond, q.Values, err = ParseSQL(cond, 0); err != nil {
		return
	}
	return q, nil
}
