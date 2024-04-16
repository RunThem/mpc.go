## mpc.go
 - ?: 零次或一次
 - *: 零次或多次
 - +: 一次或多次


```go
         
code := `
    factor = /[0-9]*/ | '(' $expr ')'       ;
    term   = $factor (('*' | '/') $factor)* ;
    expr   = $term (('+' | '-') $term)      ;
`

```
