package coroutine

import "runtime"

type Yield func(...interface{}) ([]interface{}, bool)
type Kill  func()

func Create(f func(Yield, ...interface{}) []interface{}) (Yield, Kill) {
        var signal int = 1
        var value  []interface{}
        c0 := make(chan []interface{})
        c1 := make(chan []interface{})

        go func() {
                value = <- c0
                if signal == 0 {
                        return
                }
                ret := f(func(arg...interface{}) ([]interface{}, bool) {
                        signal = 1
                        c1 <- arg
                        value = <- c0
                        if signal == 0 {
                                runtime.Goexit()
                        }
                        return value, true
                }, value...)
                signal = 0
                c1 <- ret
        }()

        y := func(arg ...interface{}) ([]interface{}, bool) {
                if signal == 0 {
                        return nil, false
                }
                signal = 1
                c0 <- arg
                value = <- c1
                return value, signal == 1
        }
        k := func() {
                if signal == 0 {
                        return
                }
                signal = 0
                c0 <- nil
        }
        return y, k
}