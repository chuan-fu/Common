package sort

const (
	two = 2
)

func Deduplice(listInter interface{}) interface{} {
	switch list := listInter.(type) {
	case []int64:
		return DedupliceInt64(list)
	case []int:
		return DedupliceInt(list)
	case []float64:
		return DedupliceFloat64(list)
	case []string:
		return DedupliceString(list)
	default:
		return nil
	}
}

func DedupliceInt64(list []int64) []int64 {
	return list[:DedupliceInt64Index(list)]
}

func DedupliceInt64Index(list []int64) int {
	num := len(list)
	if num < two {
		return num
	}
	Int64(list)

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DedupliceInt(list []int) []int {
	return list[:DedupliceIntIndex(list)]
}

func DedupliceIntIndex(list []int) int {
	num := len(list)
	if num < two {
		return num
	}
	Int(list)

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DedupliceFloat64(list []float64) []float64 {
	return list[:DedupliceFloat64Index(list)]
}

func DedupliceFloat64Index(list []float64) int {
	num := len(list)
	if num < two {
		return num
	}
	Float64(list)

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}

func DedupliceString(list []string) []string {
	return list[:DedupliceStringIndex(list)]
}

func DedupliceStringIndex(list []string) int {
	num := len(list)
	if num < two {
		return num
	}
	String(list)

	index := 1
	for i := 1; i < num; i++ {
		if list[i-1] != list[i] {
			list[index] = list[i]
			index++
		}
	}
	return index
}
