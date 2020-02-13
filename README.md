# continuously-stare-at-images
Run this and see what images your cluster is using or planning on using.

It starts a watcher and prints out the names of images used for all containers in every pod (including the init and ephemeral containers).

It is also a small example of using the kubernetes Go library.

```sh
go run main.go
```

You can the `KUBERNETES_CONFIG_PATH` if you want to not use the default macOS config file.
