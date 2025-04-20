# map

map is key-value pairs, like dict in python, and unorderedmap in c++
key shoul be support `==` and `!=`

such as

```go
const (
    StateIdle = iota
    StateNormal
    StateBusy
    StateCrazy
)

//keys is string, value is string
var params map[string]string

//keys is int(enum), value is string
var states map[int]string{
    StateIdle: "idel",
    StateNormal: "normal",
    StateBusy: "busy",
    StateCrazy: "crazy"
}

const PARAM_SIZE = 100
params: = make(map[string]string, PARAM_SIZE)
```


map 的内部实现是 hashtable, 不要依赖于 map 遍历顺序, 因为迭代器的起始位置是随机的
map 不是线程安全的, 线程安全的 map 类是 sync.Map

