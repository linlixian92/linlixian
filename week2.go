package main

import (
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
)

var ErrNoRows = errors.New("no row found")

func main() {
	err := service()
	if xerrors.Is(err, ErrNoRows) {
		fmt.Println("处理。。")
	}
	fmt.Printf("origin err: %T %v\n", xerrors.Cause(err), xerrors.Cause(err)) // 打根因类型，根因信息
	fmt.Printf("main: %+v\n", err)                                            // 打印堆栈信息
}

func service() error {
	err := QueryRows()
	return xerrors.WithMessagef(err, "service failed") // 只包装根因
}

func queryDb() error {
	err := QueryRows()
	return xerrors.Wrapf(err, "queryDb failed") // 包装根因和堆栈
}

func QueryRows() error {
	return ErrNoRows
}
