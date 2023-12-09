package main

import (
	"errors"
	"fmt"
	"os"
)

func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("in fileChecker: %w", err)
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

type MyErr struct {
	Codes int
}

func (e *MyErr) Error() string {
	return fmt.Sprintf("codes: %v", e.Codes)
}

//func (e *MyErr) Is(target error) bool {
//	if me2, ok := target.(*MyErr); ok {
//		return reflect.DeepEqual(&e, &me2)
//	}
//	return false
//}

type ResourceErr struct {
	Resource string
	Code     int
}

func (re ResourceErr) Error() string {
	return fmt.Sprintf("%s: %d", re.Resource, re.Code)
}

//func (re ResourceErr) Is(target error) bool {
//	if other, ok := target.(ResourceErr); ok {
//		ignoreResource := other.Resource == ""
//		ignoreCode := other.Code == 0
//		matchResource := other.Resource == re.Resource
//		matchCode := other.Code == re.Code
//		return matchResource && matchCode ||
//			matchResource && ignoreCode ||
//			ignoreResource && matchCode
//	}
//	return false
//}

type OriginalErr struct {
	err error
}

type ErrorCode uint64

const (
	Zero ErrorCode = iota
	One
)

func (code ErrorCode) Error() string {
	return []string{
		"Error: Zero",
		"Error: One",
	}[code]
}

func (e OriginalErr) Error() string {
	return e.err.Error()
}

func main() {
	err := fileChecker("not_here.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("That file doesn't exist")
		}
	}

	myErr := &MyErr{
		Codes: 1,
	}
	//myErr2 := myErr
	myErr2 := &MyErr{
		Codes: 1,
	}
	if errors.Is(myErr2, myErr) {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}

	err = OriginalErr{
		err: errors.New("1"),
	}
	var target OriginalErr
	as := errors.As(err, &target)
	fmt.Printf("as = %v\n", as)

	var code ErrorCode
	as = errors.As(Zero, &code)
	fmt.Printf("as = %v\n", as)
}
