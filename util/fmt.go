package util

import (
	"bytes"
	"fmt"
	"strconv"
)

type fmtGrain struct {
	i        int
	str      string
	spaceNum int
}

type FmtList struct {
	list []fmtGrain
	i    int
}

func NewFmtList() *FmtList {
	return &FmtList{
		i:    1,
		list: make([]fmtGrain, 0),
	}
}

func (p *FmtList) Add(str string) (i int) {
	p.list = append(p.list, fmtGrain{
		i:   p.i,
		str: str,
	})
	i = p.i
	p.i++
	return
}

func (p *FmtList) InitSpace() {
	length := p.GetLen(p.i - 1)
	for k := range p.list {
		v := &p.list[k]
		v.spaceNum = length - p.GetLen(v.i)
	}
}

func (p *FmtList) String() string {
	p.InitSpace()
	sBy := bytes.Buffer{}
	for k := range p.list {
		v := &p.list[k]
		sBy.WriteString(fmt.Sprintf("%d)", v.i))
		for i := 0; i <= v.spaceNum; i++ {
			sBy.WriteByte(' ')
		}
		sBy.WriteString(v.str)
		sBy.WriteString("\n")
	}
	return sBy.String()
}

func (p *FmtList) GetLen(i int) int {
	return len(strconv.Itoa(i))
}

func (p *FmtList) Len() int {
	return p.i - 1
}
