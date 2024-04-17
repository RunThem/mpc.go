## mpc.go

- ?: 零次或一次
- *: 零次或多次
- +: 一次或多次

```go
package main

const code = `
    $factor = R<[0-9]*>R | '(' $expr ')'       ;
    $term   = $factor (('*' | '/') $factor)* ;
    $expr   = $term (('+' | '-') $term)      ;
`

```
