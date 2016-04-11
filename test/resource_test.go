package test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/icsnju/apt-mesos/mesosproto"
	schedule "github.com/icsnju/apt-mesos/scheduler/impl"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRangeAdd(t *testing.T) {
	Convey("range add", t, func() {
		rangesA := &mesosproto.Value_Ranges{}
		range1 := &mesosproto.Value_Range{
			Begin: proto.Uint64(10),
			End:   proto.Uint64(15),
		}
		range2 := &mesosproto.Value_Range{
			Begin: proto.Uint64(21),
			End:   proto.Uint64(25),
		}
		rangesA.Range = append(rangesA.Range, range1)
		rangesA.Range = append(rangesA.Range, range2)

		rangesB := &mesosproto.Value_Ranges{}
		range3 := &mesosproto.Value_Range{
			Begin: proto.Uint64(16),
			End:   proto.Uint64(19),
		}
		range4 := &mesosproto.Value_Range{
			Begin: proto.Uint64(26),
			End:   proto.Uint64(28),
		}
		rangesA.Range = append(rangesA.Range, range3)
		rangesA.Range = append(rangesA.Range, range4)

		rangeC := schedule.RangeAdd(rangesA, rangesB)
		So(rangeC.Range[0].GetBegin(), ShouldEqual, 10)
		So(rangeC.Range[1].GetEnd(), ShouldEqual, 28)
	})
}
