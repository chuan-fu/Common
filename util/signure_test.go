package util

import (
	"fmt"
	"strings"
	"testing"
)

// https://api.xiaoyisz.com/qiehuang/ga/user/info?
// timestamp=1663773382338&nonce=FxtPykmTmSFFCJzQ&signature=7FD9CA36FDE1ECDFA02C7B90340CBD0D

// https://api.xiaoyisz.com/qiehuang/ga/user/task/list?
// timestamp=1663773382338&nonce=RjRYHYMKQ352kx6m&signature=0E1FB04FFF08F2FEBF92113CFA55016D
func TestMD5(t *testing.T) {
	s := `nonce=FxtPykmTmSFFCJzQ&timestamp=1663773382338`
	fmt.Println(strings.ToUpper(Md5(s)))
}
