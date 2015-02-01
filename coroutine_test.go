package coroutine

import (
	"testing"
)

func TestCoroutine(t *testing.T) {
        yield, _ := Create(func (yield Yield, arg ...interface{}) []interface{} {
                if len(arg) != 3 {
                        t.Errorf("Arg counts error. Want 3. Got %q.", len(arg))
                }

                if arg[0].(int) != 123 || arg[1].(string) != "hello" || *arg[2].(*string) != "world" {
                        t.Errorf("Value error. Want 123, hello, world. Got %q, %q, %q.",
                                arg[0].(int), arg[1].(string), *arg[2].(*string))
                }

                yield(456)

                return []interface{}{"done"}
        })

        world := "world"
        v, ok := yield(123, "hello", &world)
        if !ok {
                t.Error("Stat error. Want true, Got false")
        }
        if v[0].(int) != 456 {
                t.Errorf("Value error. Want 456. Got %q.", v[0].(int))
        }

        v, ok = yield()
        if ok {
                t.Error("Stat error. Want false, Got true")
        }
        if v[0].(string) != "done" {
                t.Errorf("Value error. Want done. Got %q.", v[0].(string))
        }
}

func TestKillCoroutine(t *testing.T) {
        yield, kill := Create(func (yield Yield, arg ...interface{}) []interface{} {
                yield(456)
                return []interface{}{"done"}
        })


        _, ok := yield(123, "hello")
        if !ok {
                t.Error("Stat error. Want true, Got false")
        }

        kill()

        _, ok = yield(789)

        if ok {
                t.Error("Stat error. Want false, Got true")
        }
}