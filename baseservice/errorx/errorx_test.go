package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorX(t *testing.T) {
	fmt.Println(ToErrorX(&errorX{code: 1, error: errors.New("A")}).Error())
	fmt.Println(WithCode(2, errors.New("A")).Error())
	fmt.Println(ToErrorX(errors.New("A")).Error())
	fmt.Println(New("AA").Error())
	fmt.Println(Newf("A:%d:C", 1).Error())
	fmt.Println(News("AA", "BB", "CC").Error())
	fmt.Println(Wrap(errors.New("AAA"), "key").Error())
	fmt.Println(Wrapf(errors.New("AAA"), "%d %d", 1, 2).Error())
}
