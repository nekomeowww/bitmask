# bitmask

[![build](https://github.com/nekomeowww/bitmask/actions/workflows/build.yaml/badge.svg)](https://github.com/nekomeowww/bitmask/actions/workflows/build.yaml) [![unittest](https://github.com/nekomeowww/bitmask/actions/workflows/test.yml/badge.svg)](https://github.com/nekomeowww/bitmask/actions/workflows/test.yml) [![](https://goreportcard.com/badge/github.com/nekomeowww/bitmask)](https://goreportcard.com/report/github.com/nekomeowww/bitmask)

A golang library to manipulating bitmasks with marshal/unmarshal struct support.

Features:

- Manually manipulating bitmasks with easy to use functions
- Marshal/Unmarshal bitmasks into or from structs
- Support all int, uint types, and bool in structs

## Usage

Important notice: **The first bit should be 1, not 0.**

Use `Set(int)` and `Unset(int)` to set and unset bits.   
Use `IsSet(int) bool` to determine if a bit is set.   
   
Use `Marshal(interface{})` to marshal a bitmask into a struct.   
Use `Unmarshal(interface{})` to unmarshal a bitmask from a struct.

### Manually

```go
package main

import "github.com/nekomeowww/bitmask"

func main() {
    b := bitmask.New(3)
    b.Set(1)
    b.Set(2)

    b.IsSet(1) // true
    b.IsSet(2) // true
    b.IsSet(3) // false

    b.Unset(1)
    b.IsSet(1) // false
}
```

### Marshal

Use the tag `bitmask` to specify the bitmask field.

#### Int

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
   A int `bitmask:"1"`
   B int `bitmask:"2"`
}

func main() {
    v := tS1{A: 1, B: 1}
    b, err := Marshal(v)
    if err != nil {
        log.Fatal(err)
    }

    b.IsSet(1) // true
    b.IsSet(2) // true
}
```

#### Uint

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
   A uint `bitmask:"1"`
   B uint `bitmask:"2"`
}

func main() {
    v := tS1{A: 1, B: 1}
    b, err := Marshal(v)
    if err != nil {
        log.Fatal(err)
    }

    b.IsSet(1) // true
    b.IsSet(2) // true
}
```

#### Bool

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
   A bool `bitmask:"1"`
   B bool `bitmask:"2"`
}

func main() {
    v := tS1{A: true, B: true}
    b, err := Marshal(v)
    if err != nil {
        log.Fatal(err)
    }

    b.IsSet(1) // true
    b.IsSet(2) // true
}
```

### Unmarshal

Use the tag `bitmask` to specify the bitmask field.

#### Int

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
    A int `bitmask:"1"`
    B int `bitmask:"2"`
}

func main() {
    var v tS1

    b := New(3)
    b.IsSet(1) // true
    b.IsSet(2) // true

    err := Unmarshal(bitmask, &v)
    require.NoError(err)

    v.A // 1
    v.B // 1
}
```

#### Uint

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
    A uint `bitmask:"1"`
    B uint `bitmask:"2"`
}

func main() {
    var v tS1

    b := New(3)
    b.IsSet(1) // true
    b.IsSet(2) // true

    err := Unmarshal(bitmask, &v)
    require.NoError(err)

    v.A // 1
    v.B // 1
}
```

#### Bool

```go
package main

import "github.com/nekomeowww/bitmask"

type tS1 struct {
    A bool `bitmask:"1"`
    B bool `bitmask:"2"`
}

func main() {
    var v tS1

    b := New(3)
    b.IsSet(1) // true
    b.IsSet(2) // true

    err := Unmarshal(bitmask, &v)
    require.NoError(err)

    v.A // true
    v.B // true
}
```
