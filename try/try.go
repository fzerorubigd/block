// Package catch is a simple try catch like (not exactly same) mechanism on golang error system
// the idea is to handle errors in a chain of function call, the goal is filter error and call
// only functions with correct signature.
//
//  sample :
//
//
//	try.New(
//		errors.New("string"), // The first catch is ignored
//	).Catch(func(e *os.PathError) error {
//		fmt.Println("os path error", e.Error())
//		return nil // Suppress error here
//	}).Catch(func(e error) error {
//		fmt.Println("string error", e)
//		return e // Pass the error to next level
//	})
//
//
//  this is not exactly a try/catch mechanism.
//
package try

import "reflect"

var errSample = reflect.TypeOf((*error)(nil)).Elem()

// Interface is the catch interface
type Interface interface {
	// Catch accept a function and call it if its matched with the signature
	Catch(f interface{}) Interface
	// Error return the final error in system. if a catcher change it to new function the
	// result is new error not the first err
	Error() error
}

// Block is the main structure for this package
type block struct {
	err error
}

func (b *block) match(err error, t reflect.Type) bool {
	return reflect.TypeOf(err).AssignableTo(t.In(0))
}

func (b *block) call(err error, v reflect.Value) error {
	args := []reflect.Value{reflect.ValueOf(err)}
	res := v.Call(args)
	if res[0].IsNil() {
		return nil
	}

	return res[0].Interface().(error)
}

func (b *block) validate(fnType reflect.Type) {
	if fnType.Kind() != reflect.Func {
		panic("must be a func")
	}

	if fnType.NumIn() != 1 || !fnType.In(0).Implements(errSample) {
		panic("must get exactly one argument and argument must implement error interface")
	}

	if fnType.NumOut() != 1 || !fnType.Out(0).Implements(errSample) {
		panic("must have exactly one result and result must implement error interface")
	}
}

// Catch register a catch block
func (b *block) Catch(f interface{}) Interface {
	if b.err != nil {
		fnType := reflect.TypeOf(f)
		b.validate(fnType)
		if b.match(b.err, fnType) {
			b.err = b.call(b.err, reflect.ValueOf(f))
		}
	}
	return b
}

func (b *block) Error() error {
	return b.err
}

func New(err error) Interface {
	return &block{err: err}
}
