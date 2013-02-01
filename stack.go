package normalize

import (
	"errors"
	"fmt"
)

type StackT struct {
	Data []string
	Size uint
	Top  uint
}

func NewStack(max uint) *StackT {
	s := make([]string, max)
	return &StackT{Data: s, Size: max, Top: 0}
}

func (s *StackT) Push(pushed string) error {
	n := s.Top
	if n >= s.Size-1 {
		return errors.New("Stack overflow")
	}
	s.Top++
	s.Data[n] = pushed
	return nil
}

func (s *StackT) Pop() (string, error) {
	n := s.Top
	if n == 0 {
		return string(""), errors.New("Stack underflow")
	}
	top := s.Data[n-1]
	s.Top--
	return top, nil
}

func (s *StackT) Print() {
	n := s.Top
	fmt.Println("Cap:", s.Size, "Size:", n)
	var i uint
	for i = 0; i < n; i++ {
		fmt.Printf("\t%d:\t%s\n", i, s.Data[i])
	}
}
