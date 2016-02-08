package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fzerorubigd/block/catch"
	"github.com/fzerorubigd/block/try"
)

func main() {
	// try
	try.New(&os.PathError{
		Op:   "test",
		Path: "test",
		Err:  errors.New("test"),
	}).Catch(func(e *os.PathError) error {
		fmt.Println("os", e.Error())
		return nil
	}).Catch(func(e error) error {
		fmt.Println("err", e)
		return e
	})

	try.New(
		errors.New("string"),
	).Catch(func(e *os.PathError) error {
		fmt.Println("os path error", e.Error())
		return nil
	}).Catch(func(e error) error {
		fmt.Println("string error", e)
		return e
	})

	//

	t := catch.New().Catch(
		func(e *os.PathError) error {
			fmt.Println("os", e.Error())
			return nil
		},
	).Finally(
		func(e error) {
			fmt.Println("finally call ", e.Error())
		},
	).Catch(
		func(e error) error {
			fmt.Println("simple error", e.Error())
			return nil
		},
	)

	fmt.Println(t.Try(errors.New("string")))
	fmt.Println(t.Try(&os.PathError{
		Op:   "test",
		Path: "test",
		Err:  errors.New("test")}))
}
