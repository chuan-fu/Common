package retryx

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestRetry(t *testing.T) {
	fmt.Println(Retry(context.Background(), func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("err1")
	}))
	fmt.Println("=====")
	fmt.Println(Retry(context.Background(), func(ctx context.Context) (interface{}, error) {
		return "1", nil
	}))
	fmt.Println("=====")
	fmt.Println(Retry(context.Background(), func(ctx context.Context) (interface{}, error) {
		fmt.Println("runFunc")
		return "1", nil
	}, WithIsRetryFunc(func(i interface{}) bool {
		is, _ := i.(string)
		return is == "1"
	})))
	fmt.Println("=====")
	fmt.Println(Retry(context.Background(), func(ctx context.Context) (interface{}, error) {
		fmt.Println("runFunc")
		return "1", nil
	}, WithIsRetryFunc(func(i interface{}) bool {
		is, _ := i.(string)
		return !(is == "1")
	})))
}

func TestRetryFunc(t *testing.T) {
	fmt.Println(RetryFunc(context.Background(), 3, func(ctx context.Context) (data interface{}, isRetry bool, err error) { return "22", false, nil }))
	fmt.Println(RetryFunc(context.Background(), 3, func(ctx context.Context) (data interface{}, isRetry bool, err error) { return "22", false, nil }))
}
