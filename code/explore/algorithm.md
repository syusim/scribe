ExploreExpr(e) is a function that returns a set of expressions equivalent to e.

Exploration algorithm A:

Let es be the set of expressions in E currently.
For each e in es, explore the children of e and append ExploreExpr(e) to E.
If no expressions were added to E this way, halt.
Otherwise, let es be the set of expressions added to E in this iteration and repeat.

Exploration algorithm {B,C}:

Let es be a {stack,queue} initially containing each expression in E.
Until es is empty
* pop an element e off of it
* explore each child of e
* if e is present in E already, loop
* and e to E
* insert into es ExploreExpr(e).

These only differ in the order in which expressions are generated.
Does this have any material impact?
