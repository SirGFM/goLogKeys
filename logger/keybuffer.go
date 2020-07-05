package logger

type keyBuffer struct {
    size int
    head int
    tail int
    keys []Key
    states []KeyState
}

func newKeyBuffer(size int) *keyBuffer {
    return &keyBuffer {
        size: size,
        head: 0,
        tail: 0,
        keys: make([]Key, size),
        states: make([]KeyState, size),
    }
}

func (kb *keyBuffer) advance_tail() {
    kb.tail++
    if kb.tail >= kb.size {
        kb.tail = 0
    }
}

func (kb *keyBuffer) push(k Key, s KeyState) {
    kb.keys[kb.head] = k
    kb.states[kb.head] = s
    kb.head++
    if kb.head >= kb.size {
        kb.head = 0
    }
    if kb.head == kb.tail {
        kb.advance_tail()
    }
}

func (kb *keyBuffer) pop() (k Key, s KeyState) {
    if kb.head == kb.tail {
        return
    }

    k = kb.keys[kb.tail]
    s = kb.states[kb.tail]
    kb.advance_tail()

    return
}

func (kb *keyBuffer) count_continuous() int {
    if kb.head > kb.tail {
        return kb.head - kb.tail
    } else if kb.head < kb.tail {
        return kb.size - kb.tail
    } else {
        return 0
    }
}

func (kb *keyBuffer) count() int {
    num := kb.count_continuous()
    if kb.head < kb.tail {
        return num + kb.head + 1
    } else {
        return num
    }
}

func (kb *keyBuffer) fast_pop(k []Key, s []KeyState) {
    if len(k) != len(s) || len(k) > kb.count_continuous() {
        panic("logger: Can't use fast_pop() for the requested amount")
    }

    from := kb.tail
    to := kb.tail + len(k)
    copy(k, kb.keys[from:to])
    copy(s, kb.states[from:to])

    if to >= kb.size {
        kb.tail = 0
    } else {
        kb.tail = to
    }
}
