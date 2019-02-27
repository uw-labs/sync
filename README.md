# rungroup
Alternative to [error group](https://godoc.org/golang.org/x/sync/errgroup) that stops as soon as any function started in the group terminates. 
It also allows users to run functions asynchronously, so that the call to wait won't wait for them to terminate, but the group will still stop
(the underlying context will be cancelled) as soon as they terminate.