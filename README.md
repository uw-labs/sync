# sync
Extension of types from [golang.org/x/sync](https://godoc.org/golang.org/x/sync).

## rungroup 
Alternative to error group that stops (i.e. cancels the underlying context) as soon as any
function started in the group terminates. For this to work it can only be created with a context.

## gogroup 
Another alternative to error group that only waits until any single function started in the group terminates.
Like rungroup it can only be created with a context and this context is cancelled as soon as any function
started in the group terminates.

NOTE: calling wait without starting any goroutine with the `Go` method will block until the parent context is canceled.
