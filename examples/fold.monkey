let foldl = fn(initial, f, xs) {
    if (len(xs) > 0) {
        return foldl(f(initial, head(xs)), f, tail(xs))
    }
    return initial
}

let sum = fn(xs) { foldl(0, fn(acc, x) { acc + x }, xs) }
let prod = fn(xs) { foldl(1, fn(acc, x) { acc * x }, xs) }
let join = fn(xs, sep) {
    if (len(xs) < 1) {
        return ""
    }
    return foldl(head(xs), fn(acc, x) { acc + sep + x }, tail(xs))
}

print(sum([1, 2, 3, 4]))
print(prod([1, 2, 3, 4]))
print(join(["foo", "bar", "baz"], ", "))
