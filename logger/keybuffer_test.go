package logger

import (
    crypto_rand "crypto/rand"
    "math/rand"
    "testing"
)

func seed() {
    var data []byte = make([]byte, 8)

    n, err := crypto_rand.Read(data)
    if n != len(data) || err != nil {
        panic("Failed to seed the RNG")
    }

    var val int64
    val  = int64(data[0] & 0xff)
    val |= int64(data[1] & 0xff) << 8
    val |= int64(data[2] & 0xff) << 16
    val |= int64(data[3] & 0xff) << 24
    val |= int64(data[4] & 0xff) << 32
    val |= int64(data[5] & 0xff) << 40
    val |= int64(data[6] & 0xff) << 48
    val |= int64(data[7] & 0xff) << 56
    rand.Seed(val)
}

func TestKeyBuffer(t *testing.T) {
    var in []Key
    var out_k []Key
    var out_s []KeyState

    seed()

    in = make([]Key, 12)
    for i := range in {
        in[i] = Key(rand.Intn(int(KeyCount)))
    }

    // Basic test: fill buffer and pop everything
    buf := newKeyBuffer(len(in) + 1)
    for _, v := range in {
        buf.push(v, Released)
    }
    for _, v := range in {
        k, _ := buf.pop()
        if want, got := v, k; want != got {
            t.Errorf("Invalid data. Want: %+v, Got: %+v", want, got)
        }
    }

    // Check if a few entry are overriden by the circular loop
    skip := 4
    buf = newKeyBuffer(len(in) + 1 - skip)
    for _, v := range in {
        buf.push(v, Released)
    }
    for i := range in {
        if skip > 0 {
            skip--
            continue
        }

        k, _ := buf.pop()
        if want, got := in[i], k; want != got {
            t.Errorf("Invalid data. Want: %+v, Got: %+v", want, got)
        }
    }

    // Check if the fast interface works
    out_k = make([]Key, len(in))
    out_s = make([]KeyState, len(in))
    buf = newKeyBuffer(len(in) + 1)
    for _, v := range in {
        buf.push(v, Released)
    }
    buf.fast_pop(out_k, out_s)
    for i := range in {
        if want, got := in[i], out_k[i]; want != got {
            t.Errorf("Invalid data. Want: %+v, Got: %+v", want, got)
        }
    }

    // Check if the fast interface works after overriding entries
    skip = 4
    buf = newKeyBuffer(len(in) + 1 - skip)
    for _, v := range in {
        buf.push(v, Released)
    }

    count := buf.count_continuous()
    if count == len(in) || count == 0 {
        t.Error("Continuout buffer shouldn't be equal to input")
    }

    test := in[skip:skip+count]
    out_k = make([]Key, count)
    out_s = make([]KeyState, count)
    buf.fast_pop(out_k, out_s)
    for i := range test {
        if want, got := test[i], out_k[i]; want != got {
            t.Errorf("Invalid data. Want: %+v, Got: %+v", want, got)
        }
    }

    count = buf.count_continuous()
    test = in[skip + len(out_k):]
    if count != len(test) || count == 0 {
        t.Error("Invalid length")
    }

    out_k = make([]Key, count)
    out_s = make([]KeyState, count)
    buf.fast_pop(out_k, out_s)
    for i := range test {
        if want, got := test[i], out_k[i]; want != got {
            t.Errorf("Invalid data. Want: %+v, Got: %+v", want, got)
        }
    }

    count = buf.count_continuous()
    if count != 0 {
        t.Error("Expected keybuf to be empty!")
    }
}
