package objects

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestObjects(t *testing.T) {
	Reset()

	var (
		id = NewObjectId("100")
	)

	Convey("Check", t, func() {
		So(id, ShouldEqual, 1)
		obj := id.Get().(string)
		So(obj, ShouldEqual, "100")
	})

	Convey("Free", t, func() {
		So(len(manager.freeObjs), ShouldEqual, 0)
		id.Free()
		So(len(manager.freeObjs), ShouldEqual, 1)
	})

	Convey("New Again id should == 1", t, func() {
		id = NewObjectId("200")

		So(id, ShouldEqual, 1)
		obj := id.Get().(string)
		So(obj, ShouldEqual, "200")

		id.Free()
	})
}

func TestParallel(t *testing.T) {
	Convey("Parallel", t, func(c C) {
		Reset()
		wait := sync.WaitGroup{}
		wait.Add(100)

		for i := 0; i < 100; i++ {
			go func(c C) {
				for j := 0; j < 100; j++ {
					o := strconv.Itoa(j + i*100)
					id := NewObjectId(o)
					time.Sleep(time.Duration(rand.Int63n(1000)))
					c.So(id.Get().(string), ShouldEqual, o)
					id.Free()
				}

				wait.Done()
			}(c)
		}

		wait.Wait()

		Convey("Check", func() {
			t.Log("ID next: ", manager.next)
			t.Log("len objs: ", len(manager.objs))
			t.Log("len freeObjs:", len(manager.freeObjs))

			for key, _ := range manager.objs {
				t.Log("Objs:", key)
			}

			So(len(manager.objs), ShouldEqual, 0)
			So(len(manager.freeObjs), ShouldBeGreaterThan, 99)
		})
	})
}
