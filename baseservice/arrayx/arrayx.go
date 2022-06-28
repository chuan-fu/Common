package arrayx

import (
	"bytes"
	"strconv"

	"github.com/chuan-fu/Common/cdefs"
)

func IsInArray(s int64, array []int64) bool {
	for _, v := range array {
		if s == v {
			return true
		}
	}
	return false
}

func IsInStringArray(key string, list []string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

func Int64Join(list []int64, sep string) string {
	if len(list) == 0 {
		return ""
	}
	b := bytes.Buffer{}
	for k, v := range list {
		if k > 0 {
			b.WriteString(sep)
		}
		b.WriteString(strconv.FormatInt(v, cdefs.BitSize10))
	}
	return b.String()
}

// 拆分array
func SplitArray(arr []int64, num int64) [][]int64 {
	if len(arr) == 0 {
		return [][]int64{}
	}
	max := int64(len(arr))
	// 判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]int64{arr}
	}
	// 获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	// 声明分割好的二维数组
	segments := make([][]int64, 0)
	// 声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func ConvertToIntArray(arr []string) ([]int, bool) {
	result := make([]int, 0)
	for _, i := range arr {
		res, err := strconv.Atoi(i)
		if err != nil {
			return result, false
		}
		result = append(result, res)
	}
	return result, true
}
