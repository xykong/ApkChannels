package sign

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIntegerStuff(t *testing.T) {
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}

func TestSeekEocdWithCommentInBuffer(t *testing.T) {
	Convey("The value should be greater by one", t, func() {
		So(2, ShouldEqual, 2)
	})
}
