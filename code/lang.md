Design Ideas:

* Variadic stuff is kind of messy in a language where you have a
  strong distinction between relexprs and scalar exprs, so opt for
  lists.
* I think what makes this hard is that so much stuff gets added to
  every scope invisibly, and sexps force you to make things very
  explicit.

`create-table` expressions must contain their data inline.
```
(create-table foo
  [[a int] [b int] [c int]]
  [[1 2 3]
   [4 5 6]])

(create-table bar
  [[x int]
   [y int]
   [z int]]
  [[10 20 30]
   [40 50 60]])
```

To read a table you just give the name directly.
```
foo
```

A relational expression exposes a *scope* which is simply
a list of Columns.
```
scope(foo)
=>
[[foo a int] [foo b int] [foo c int]]

scope((join foo bar))
=>
[[foo a int] [foo b int] [foo c int]
 [bar x int] [bar y int] [bar z int]]
```

The table of a scope can be changed with the `as` operator.
```
scope((as foo bar))
=>
[[bar a int] [bar b int] [bar c int]]
```

`as` also allows us to rename columns:
```
scope((as foo bar [u _ w]))
=>
[[bar u int] [bar b int] [bar w int]]
```

`project` allows us to render scalar expressions.
```
(project foo [a b])

(project foo [(+ a b)])

(project foo [(as (+ a b) sum)])
```

`select` filters rows from an expression.
```
(select foo
  (and (> a 1) (= c 3)))
```

```
(join foo bar (= a x))

(join foo (as foo bar) (= foo.a bar.a))

(left-join foo (as foo bar) (= foo.a bar.a))
```

```
(values
  [[1 2 3]
   [4 5 6]])
```

There's a special operator that is only valid at the root, called `order-by`.
`order-by` allows you to request a particular order on the data you will receive.
```
(order-by (select foo (= x 3))
  [y])
```
