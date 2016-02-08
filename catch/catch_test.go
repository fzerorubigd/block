package catch

import (
	"errors"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MyError struct {
	Data string
}

func (me *MyError) Error() string {
	return me.Data
}

func TestTryCatch(t *testing.T) {
	Convey("Test try/catch", t, func() {
		var tmp string
		var finCall int
		last := New().Catch(
			func(err *os.PathError) error {
				panic("must not call this")
				return nil
			},
		).Finally(
			func(err error) {
				finCall++
			},
		).Catch(
			func(err error) error {
				tmp = err.Error()
				return nil
			},
		).Catch(
			func(err error) error {
				panic("must not call this")
				return nil
			},
		).Finally(
			func(err error) {
				finCall += 2
			},
		)
		So(tmp, ShouldEqual, "")
		So(last.Try(errors.New("string")), ShouldBeNil)
		So(tmp, ShouldEqual, "string")
		So(finCall, ShouldEqual, 3)
	})

	Convey("Test try/catch chain", t, func() {
		var tmp string
		So(New().Catch(
			func(err *os.PathError) error {
				panic("must not call this")
				return nil
			},
		).Catch(
			func(err error) error {
				tmp = err.Error()
				return err
			},
		).Catch(
			func(err error) error {
				tmp += err.Error()
				return nil
			},
		).Try(errors.New("string")), ShouldBeNil)

		So(tmp, ShouldEqual, "stringstring")
	})

	Convey("Test try/catch chain multiple", t, func() {
		var call int
		So(New().Catch(
			func(err *os.PathError) error {
				call++
				return err.Err
			},
		).Catch(
			func(err error) error {
				call++
				return err
			},
		).Catch(
			func(err error) error {
				call++
				return err
			},
		).Try(&os.PathError{Err: errors.New("string")}), ShouldNotBeNil)

		So(call, ShouldEqual, 3)
	})

	Convey("Test try/catch chain break", t, func() {
		var call int
		New().Catch(
			func(err *os.PathError) error {
				call++
				return err.Err
			},
		).Catch(
			func(err error) error {
				call++
				return nil
			},
		).Catch(
			func(err error) error {
				call++
				return nil
			},
		).Try(&os.PathError{Err: errors.New("string")})

		So(call, ShouldEqual, 2)
	})

	Convey("Test invalid input", t, func() {
		So(func() {
			New().Catch(1)
		}, ShouldPanic)

		So(func() {
			New().Catch(
				func(err int) error {
					return nil
				},
			)
		}, ShouldPanic)

		So(func() {
			New().Catch(
				func(err error, extra int) error {
					return nil
				},
			)
		}, ShouldPanic)

		So(func() {
			New().Catch(
				func(err error) int {
					return 0
				},
			)
		}, ShouldPanic)
	})
}
