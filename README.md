## dummy repo for testing out [cloudevents.io][ce]

### NATs pub/example

Install and start NATs:
```
docker pull nats:latest
docker run -p 4222:4222 -ti nats:latest
```

Start `party` api:
```
go run party/main.go
```

Start the `jokeprovider`:
```
go run jokeprovider/main.go
```

Curl the API:

```
curl localhost:8090/enrol?target=jokeprovider
```

The API will publish an event to the `jokeprovider` subject. The `jokeprovider`
will be subscribed to that subject. The event will have the `action` instruction
which is based on the original URL. The provider will perform the `action`.

### sender per provider example

TODO

[ce]: https://cloudevents.io/
