What are all the things that need to happen for a new relational expr?

1. Need to add a new struct for it in memo/rel.go with appropriate methods.
2. Need to add a method for in intern (this can maybe be combined with (1)
3. Need to augment memo/walk.go
