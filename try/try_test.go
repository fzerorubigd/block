package try

import (
	"errors"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTryCatch(t *testing.T) {
	Convey("Test try/catch", t, func() {
		var tmp string
		last := New(errors.New("string")).Catch(
			func(err *os.PathError) error {
				So(false, ShouldBeTrue)
				return nil
			},
		).Catch(
			func(err error) error {
				tmp = err.Error()
				return nil
			},
		).Catch(
			func(err error) error {
				So(false, ShouldBeTrue)
				return nil
			},
		).Error()

		So(tmp, ShouldEqual, "string")
		So(last, ShouldBeNil)
	})

	Convey("Test try/catch chain", t, func() {
		var tmp string
		So(New(errors.New("string")).Catch(
			func(err *os.PathError) error {
				So(false, ShouldBeTrue)
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
		).Error(), ShouldBeNil)

		So(tmp, ShouldEqual, "stringstring")
	})

	Convey("Test try/catch chain multiple", t, func() {
		var call int
		So(New(&os.PathError{Err: errors.New("string")}).Catch(
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
		).Error(), ShouldNotBeNil)

		So(call, ShouldEqual, 3)
	})

	Convey("Test try/catch chain break", t, func() {
		var call int
		New(&os.PathError{Err: errors.New("string")}).Catch(
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
		)

		So(call, ShouldEqual, 2)
	})

	Convey("Test invalid input", t, func() {
		So(func() {
			New(errors.New("string")).Catch(1)
		}, ShouldPanic)

		So(func() {
			New(errors.New("string")).Catch(
				func(err int) error {
					return nil
				},
			)
		}, ShouldPanic)

		So(func() {
			New(errors.New("string")).Catch(
				func(err error, extra int) error {
					return nil
				},
			)
		}, ShouldPanic)

		So(func() {
			New(errors.New("string")).Catch(
				func(err error) int {
					return 0
				},
			)
		}, ShouldPanic)
	})
}
