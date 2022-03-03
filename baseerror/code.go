package baseerror

const (
	NotFoundErrorCode = 404
)

const (
	NotFoundErrorMsg = "没有对应记录"
)

var errCodeMsg = map[int]string{
	NotFoundErrorCode: NotFoundErrorMsg,
}
