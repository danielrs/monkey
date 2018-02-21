let fizzbuzz = fn(n) {
    let iter = fn(i, n) {
        if (i < n) {
            let msg =
                (i % 15 == 0 && "Fizz Buzz")
                || (i % 3 == 0 && "Fizz")
                || (i % 5 == 0 && "Buzz")
                || i
            print(msg)
            iter(i + 1, n)
        }
    }
    iter(1, n)
}

fizzbuzz(100)
