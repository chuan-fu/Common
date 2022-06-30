package arrayx

/*
	排序+去重
*/

func DistinctInt64(list []int64) []int64 {
	return list[:DistinctInt64Index(list)]
}

func DistinctInt64Index(list []int64) int {
	num := len(list)
	if num < 2 {
		return num
	}
	Int64(list) // sort

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DistinctInt(list []int) []int {
	return list[:DistinctIntIndex(list)]
}

func DistinctIntIndex(list []int) int {
	num := len(list)
	if num < 2 {
		return num
	}
	Int(list) // sort

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DistinctFloat64(list []float64) []float64 {
	return list[:DistinctFloat64Index(list)]
}

func DistinctFloat64Index(list []float64) int {
	num := len(list)
	if num < 2 {
		return num
	}
	Float64(list) // sort

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DistinctString(list []string) []string {
	return list[:DistinctStringIndex(list)]
}

func DistinctStringIndex(list []string) int {
	num := len(list)
	if num < 2 {
		return num
	}
	String(list) // sort

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}
