# Block 
[![Build Status](https://travis-ci.org/fzerorubigd/block.svg?branch=master)](https://travis-ci.org/fzerorubigd/block)

A simple try/catch like mechanism for golang. My goal is not implement the exact try/catch model base 
on panic, but a simple way to handle errors with less if/else in the code, especially when you need
to test the error against multiple error type.
 
*every thing is just an experiment and api is going to change. I use this code in production, but I do not 
recommend you to do the same*

## try  
[![GoDoc](https://godoc.org/github.com/fzerorubigd/block/try?status.svg)](https://godoc.org/github.com/fzerorubigd/block/try)
try is the simple one, more like the try/catch block and try is in the begin of the block :

```go 
	try.New(
		err,
	).Catch(func(e *os.PathError) error { // if error is from this type, this function is executed 
		fmt.Println("os path error", e.Error())
		return nil  // nil means block the chain.
	}).Catch(func(e error) error {
		fmt.Println("string error", e)
		return e // means the block can continue, if there is another catch in the way
	}).Error() // return the last executed block result.  

```

this feels like the normal try/catch, but the problem is its not re-usable unless put it in a function and call
that function, and also there is no way to support finally.


## catch 
[![GoDoc](https://godoc.org/github.com/fzerorubigd/block/catch?status.svg)](https://godoc.org/github.com/fzerorubigd/block/catch)
catch is the reverse try/catch/finally. first setup all Catch and Finally blocks and the call Try. 

```go
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

	t.Try(errors.New("string"))
	
	// re-use it.
	t.Try(&os.PathError{
		Op:   "test",
		Path: "test",
		Err:  errors.New("test")},
	)
```

the finally blocks are called after all catch blocks, using defer, but there is a check and if the error is 
not match with the finally parameter, the finally is ignored, may be I change this behavior. 

## block functions

Catch functions, must accept one and only one argument support the error interface, and return exactly one 
 error supported interface.
 
if the result is nil, the chain is blocked and no other catch is executed. (but all finally are executed in reverse order) 
if the result is not nil, the next Catch is called with the new error, not the old one. 

Finally functions are different. the result is not important and the caller totally ignore the result. if the parameter is 
 the `error` type, Finally is called at any call to Try, but if the parameter is some other error type, then 
 the block is called only if the `Try()` parameter is acceptable in the finally function
 
 
Made in a boring day when I was sick. so sorry if this is not what golang needed :)