package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFCFSQueue(t *testing.T) {
	Convey("fcfs queue sort", t, func() {
		var arr []*Task
		arr = append(arr, &Task{
			ID:          "1",
			CreatedTime: 65,
		})
		arr = append(arr, &Task{
			ID:          "2",
			CreatedTime: 75,
		})
		arr = append(arr, &Task{
			ID:          "3",
			CreatedTime: 45,
		})
		queue := NewFCFSQueue(arr)
		So(queue[0].ID, ShouldEqual, "3")
		So(queue[1].ID, ShouldEqual, "1")
		So(queue[2].ID, ShouldEqual, "2")
	})
}
