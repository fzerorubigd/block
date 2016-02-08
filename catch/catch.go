// Package catch is the reverse try/catch/finally. first setup all Catch and Finally blocks and the call Try.
//
//	t := catch.New().Catch(
//		func(e *os.PathError) error {
//			fmt.Println("os", e.Error())
//			return nil
//		},
//	).Finally(
//		func(e error) {
//			fmt.Println("finally call ", e.Error())
//		},
//	).Catch(
//		func(e error) error {
//			fmt.Println("simple error", e.Error())
//			return nil
//		},
//	)
//
//	t.Try(errors.New("string"))
//
//	// re-use it.
//	t.Try(&os.PathError{
//		Op:   "test",
//		Path: "test",
//		Err:  errors.New("test")},
//	)
//
//the finally blocks are called after all catch blocks, using defer, but there is a check and if the error is
//not match with the finally parameter, the finally is ignored, may be I change this behavior.
package catch

import "reflect"

var errSample = reflect.TypeOf((*error)(nil)).Elem()

// Interface is the catch interface
type Interface interface {
	// Catch accept a function and call it if its matched with the signature
	Catch(f interface{}) Interface
	// Finally is running after all catch in reverse order
	Finally(f interface{}) Interface
	// Try must call this after adding all catch and finalies
	Try(error) error
}

type item struct {
	fnType  reflect.Type
	fnValue reflect.Value
}

// Block is the main structure for this package
type block struct {
	catch   []item
	finally []item
}

func (i item) match(err error) bool {
	return reflect.TypeOf(err).AssignableTo(i.fnType.In(0))
}

func (i item) call(err error) error {
	args := []reflect.Value{reflect.ValueOf(err)}
	res := i.fnValue.Call(args)
	if res[0].IsNil() {
		return nil
	}

	return res[0].Interface().(error)
}

func (i item) callIgnore(err error) {
	args := []reflect.Value{reflect.ValueOf(err)}
	// ignore the result
	i.fnValue.Call(args)
}

func (b *block) validate(fnType reflect.Type, ret bool) {
	if fnType.Kind() != reflect.Func {
		panic("must be a func")
	}

	if fnType.NumIn() != 1 || !fnType.In(0).Implements(errSample) {
		panic("must get exactly one argument and argument must implement error interface")
	}

	if ret {
		if fnType.NumOut() != 1 || !fnType.Out(0).Implements(errSample) {
			panic("must have exactly one result and result must implement error interface")
		}
	}
}

// Catch register a catch block
func (b *block) Catch(f interface{}) Interface {
	fnType := reflect.TypeOf(f)
	b.validate(fnType, true)
	b.catch = append(
		b.catch,
		item{
			fnType:  fnType,
			fnValue: reflect.ValueOf(f),
		},
	)
	return b
}

// Catch register a catch block
func (b *block) Finally(f interface{}) Interface {
	fnType := reflect.TypeOf(f)
	b.validate(fnType, false)
	b.finally = append(
		b.finally,
		item{
			fnType:  fnType,
			fnValue: reflect.ValueOf(f),
		},
	)
	return b
}

func (b *block) runFinally(err error) {
	for j := range b.finally {
		// some tools report problem with defer in loop. but I don't think its a problem here
		tmp := b.finally[j]
		defer func() {
			if tmp.match(err) {
				tmp.callIgnore(err)
			}
		}()
	}
}

// Try is called on error and return the error from the catch blocks
func (b *block) Try(err error) error {
	defer b.runFinally(err)
	for i := range b.catch {
		if b.catch[i].match(err) {
			err = b.catch[i].call(err)
			if err == nil {
				break
			}
		}
	}
	return err
}

// New return the
func New() Interface {
	return &block{}
}
