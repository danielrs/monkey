// Folding functions.

// foldl traverse the array left-to-right,
// using f to generate a new value every
// iteration.
let foldl = fn(initial, f, xs) {
    if (len(xs) < 1) {
        return initial
    }
    return foldl(f(initial, head(xs)), f, tail(xs))
}

// foldr is just like foldl but traverses
// the array right-to-left.
let foldr = fn(initial, f, xs) {
    if (len(xs) < 1) {
        return initial
    }
    return foldr(f(last(xs), initial), f, tail(xs))
}

// Common aggregators.

let sum = fn(xs) { foldl(0, fn(acc, x) { acc + x }, xs) }
let prod = fn(xs) { foldl(1, fn(acc, x) { acc * x }, xs) }
let map = fn(f, xs) { foldl([], fn(acc, x) { push(acc, f(x)) }, xs) }

let filter = fn(pred, xs) {
    foldl([], fn(acc, x) {
        if (pred(x)) {
            return push(acc, x)
        }
        acc
    }, xs)
}

let join = fn(xs, sep) {
    if (len(xs) < 1) {
        return ""
    }
    return foldl(head(xs), fn(acc, x) { acc + sep + x }, tail(xs))
}

// Testing.

let arr = [1, 2, 3, 4]

print(sum(arr))
print(prod(arr))
print(map(fn(x) { x * 2 }, arr))
print(filter(fn(x) { x % 2 != 0 }, arr))
print(join(["foo", "bar", "baz"], ", "))
