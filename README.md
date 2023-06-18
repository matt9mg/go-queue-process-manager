# go-queue-process-manager
Maintains a number of processes running until they are all complete

### Installation
```
go get github.com/matt9mg/go-queue-process-manager
```

### Usage Examples
##### Basic Example
```go
q := queue_process_manager.NewQueue()

for _, item := range items {
    q.AddToQueue(&queue_process_manager.QueueItem{
        QueueFunc: func(args []any) {
            myItem := args[0].(*MyItem) // converting my argument
            // do some business logic with arg
        },
        QueueFuncArgs: []any{item},
    })
}

q.ProcessQueue()
```

##### With Custom Example
```go
q := queue_process_manager.NewQueue(queue_process_manager.WithCustomMaxQueueAllowance(25))

for _, item := range items {
    q.AddToQueue(&queue_process_manager.QueueItem{
        QueueFunc: func(args []any) {
            myItem := args[0].(*MyItem) // converting my argument
            // do some business logic with arg
        },
        QueueFuncArgs: []any{item},
    })
}

q.ProcessQueue()
```

##### With Multiple Args Example
```go
q := queue_process_manager.NewQueue()

q.AddToQueue(&queue_process_manager.QueueItem{
    QueueFunc: func(args []any) {
        myItem := args[0].(*MyItem) // converting my argument
		aes256Key := args[1].(string)
        // do some business logic with arg
    },
    QueueFuncArgs: []any{item, "aes256key"},
})

q.ProcessQueue()
```