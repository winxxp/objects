package objects

import (
	"github.com/smartystreets/goconvey/convey"
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

	convey.Convey("Check", t, func() {
		convey.So(id, convey.ShouldEqual, 1)
		obj := id.Get().(string)
		convey.So(obj, convey.ShouldEqual, "100")
	})

	convey.Convey("Free", t, func() {
		convey.So(len(manager.freeObjs), convey.ShouldEqual, 0)
		id.Free()
		convey.So(len(manager.freeObjs), convey.ShouldEqual, 1)
	})

	convey.Convey("New Again id should == 1", t, func() {
		id = NewObjectId("200")

		convey.So(id, convey.ShouldEqual, 1)
		obj := id.Get().(string)
		convey.So(obj, convey.ShouldEqual, "200")

		id.Free()
	})
}

func TestParallel(t *testing.T) {
	convey.Convey("Parallel", t, func(c convey.C) {
		Reset()
		wait := sync.WaitGroup{}
		wait.Add(100)

		for i := 0; i < 100; i++ {
			go func(c convey.C) {
				for j := 0; j < 100; j++ {
					o := strconv.Itoa(j + i*100)
					id := NewObjectId(o)
					time.Sleep(time.Duration(rand.Int63n(1000)))
					c.So(id.Get().(string), convey.ShouldEqual, o)
					id.Free()
				}

				wait.Done()
			}(c)
		}

		wait.Wait()

		convey.Convey("Check", func() {
			t.Log("ID next: ", manager.next)
			t.Log("len objs: ", len(manager.objs))
			t.Log("len freeObjs:", len(manager.freeObjs))

			for key, _ := range manager.objs {
				t.Log("Objs:", key)
			}

			convey.So(len(manager.objs), convey.ShouldEqual, 0)
			convey.So(len(manager.freeObjs), convey.ShouldBeGreaterThan, 99)
		})
	})
}
