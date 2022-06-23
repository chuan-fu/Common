package arrayx

import "sort"

type i6Slice []int64

func (l i6Slice) Len() int           { return len(l) }
func (l i6Slice) Less(i, j int) bool { return l[i] < l[j] }
func (l i6Slice) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func isReverse(seq []bool) bool {
	if len(seq) > 0 && seq[0] {
		return true
	}
	return false
}

func Int64(list []int64, seq ...bool) {
	if isReverse(seq) {
		sort.Sort(sort.Reverse(i6Slice(list)))
		return
	}
	sort.Sort(i6Slice(list))
}

func Int(list []int, seq ...bool) {
	if isReverse(seq) {
		sort.Sort(sort.Reverse(sort.IntSlice(list)))
		return
	}
	sort.Ints(list)
}

func String(list []string, seq ...bool) {
	if isReverse(seq) {
		sort.Sort(sort.Reverse(sort.StringSlice(list)))
		return
	}
	sort.Strings(list)
}

func Float64(list []float64, seq ...bool) {
	if isReverse(seq) {
		sort.Sort(sort.Reverse(sort.Float64Slice(list)))
		return
	}
	sort.Float64s(list)
}
