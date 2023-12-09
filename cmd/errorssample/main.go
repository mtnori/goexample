package main

import (
	"errors"
	"fmt"
)

type MyError1Interface interface {
	Error() string
	Echo() string
}

type MyError1 struct{}

func (e *MyError1) Error() string {
	return "this is my error1"
}

func (e *MyError1) Echo() string {
	return "hello world"
}

type MyError2 struct{}

func (e *MyError2) Error() string {
	return "this is my error2"
}

func main() {
	var myError1Interface MyError1Interface

	err11 := &MyError1{}
	err12 := &MyError1{}
	err21 := &MyError2{}
	err22 := &MyError2{}

	wrappedErr11 := fmt.Errorf("wrapped %w", err11)

	fmt.Printf("error.As(err11, &myError1Interface) = %v\n", errors.As(err11, &myError1Interface))
	fmt.Printf("errors.As(err11, &err11) = %v\n", errors.As(err11, &err11))
	fmt.Printf("errors.As(err11, &err12) = %v\n", errors.As(err11, &err12))
	fmt.Printf("errors.As(err11, &err21) = %v\n", errors.As(err11, &err21))
	fmt.Printf("errors.As(err11, &err22) = %v\n", errors.As(err11, &err22))

	fmt.Println("------------------------------------------")

	fmt.Printf("errors.As(err21, &myError1Interface) = %v\n", errors.As(err21, &myError1Interface))
	fmt.Printf("errors.As(err21, &err11) = %v\n", errors.As(err21, &err11))
	fmt.Printf("errors.As(err21, &err12) = %v\n", errors.As(err21, &err12))
	fmt.Printf("errors.As(err21, &err21) = %v\n", errors.As(err21, &err21))
	fmt.Printf("errors.As(err21, &err22) = %v\n", errors.As(err21, &err22))

	fmt.Println("------------------------------------------")

	fmt.Printf("errors.As(wrappedErr11, &myError1Interface) = %v\n", errors.As(wrappedErr11, &myError1Interface))
	fmt.Printf("errors.As(wrappedErr11, &err11) = %v\n", errors.As(wrappedErr11, &err11))
	fmt.Printf("errors.As(wrappedErr11, &err12) = %v\n", errors.As(wrappedErr11, &err12))
	fmt.Printf("errors.As(wrappedErr11, &err21) = %v\n", errors.As(wrappedErr11, &err21))
	fmt.Printf("errors.As(wrappedErr11, &err22) = %v\n", errors.As(wrappedErr11, &err22))

	fmt.Println("------------------------------------------")

	err1 := errors.New("aaa")
	err2 := errors.New("bbb")
	wrappedErr1 := fmt.Errorf("wrapped %w", err1)
	wrappedWrappedErr1 := fmt.Errorf("wrapped %w", wrappedErr1)

	// 特定のエラーであるのか、特定のエラーをWrapしたエラーであるのか
	fmt.Printf("errors.Is(err1, err2) = %v\n", errors.Is(err1, err2))
	fmt.Printf("errors.Is(wrappedErr1, err1) = %v\n", errors.Is(wrappedErr1, err1))
	fmt.Printf("errors.Is(wrappedWrappedErr1, err1) = %v\n", errors.Is(wrappedWrappedErr1, err1))
	fmt.Printf("errors.Is(wrappedWrappedErr1, wrappedErr1) = %v\n", errors.Is(wrappedWrappedErr1, wrappedErr1))
}
