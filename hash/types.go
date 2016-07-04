package main

import (
	"fmt"
	"path/filepath"
)

type Size struct {
	size      int64
	humanRead bool
}

func NewSize(size int64, humanRead bool) *Size {
	return &Size{
		size:      size,
		humanRead: humanRead,
	}
}

func (s *Size) String() string {
	if s.humanRead {
		return s.HumanReadable()
	}
	return fmt.Sprint(s.size)
}

func (s *Size) HumanReadable() string {
	size := float64(s.size)
	if size < 1024 {
		return fmt.Sprintf("%.0fB", size)
	}
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%.1fK", size)
	}
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%.1fM", size)
	}
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%.1fG", size)
	}
	size /= 1024
	return fmt.Sprintf("%.1fT", size)
}

type HashValue struct {
	h     []byte
	upper bool
}

func (h HashValue) String() string {
	f := "%x"
	if h.upper {
		f = "%X"
	}
	return fmt.Sprintf(f, h.h)
}

type HashInfo struct {
	Name     string
	FullName string
	Hash     HashValue
	Size     *Size
}

func NewHashInfo(name string, size *Size, hash []byte, upper bool) *HashInfo {
	abs, _ := filepath.Abs(name)
	return &HashInfo{
		Name:     name,
		FullName: abs,
		Size:     size,
		Hash:     HashValue{hash, upper},
	}
}

// type Formater struct {
// 	*template.Template
// }

// func NewFormater(ft string) (*Formater, error) {
// 	tpl, err := template.New("").Parse(ft)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Formater{tpl}, nil
// }
