package model

//sort map by value
type PairList_StrInt []Pair_StrInt
type Pair_StrInt struct {
	Key   string
	Value int
}

func (p PairList_StrInt) Len() int {
	return len(p)
}
func (p PairList_StrInt) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}
func (p PairList_StrInt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
