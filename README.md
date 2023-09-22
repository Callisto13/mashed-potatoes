## dummy repo for testing out [cloudevents.io][ce]

I wouldn't read the code for any of this as I really wasn't trying for style or
cleanliness.

### NATs pub/example

Install and start NATs:
```
docker pull nats:latest
docker run -p 4222:4222 -ti nats:latest
```

Start `party` api:
```
go run party/main.go --protocol nats
```

Start the `jokeprovider`:
```
go run jokeprovider/nats/main.go
```

Curl the API:

```
curl localhost:8090/enrol?target=jokeprovider
```

The API will publish an event to the `jokeprovider` subject. The `jokeprovider`
will be subscribed to that subject. The event will have the `action` instruction
which is based on the original URL. The provider will perform the `action`.

The benefit to this model is that it is easy to register new Providers with the
system. The new provider simply needs to start subscribing to their own name,
and Party can then forward actions and data to the correct one without really
knowing much about it.

### http sender per provider example

Start `party` api:
```
go run party/main.go --protocol http
```

Start the `jokeprovider`:
```
go run jokeprovider/http/main.go
```

Curl the API:

```
curl localhost:8090/enrol?target=jokeprovider
```

In this case party is started with a list of registered providers. The targeted
provider will be pulled from the query and passed to the http event sender. The
event is sent to the correct provider, which executes the action.

For this party would need to read config at start up to be able to know the
correct destination of each Provider, so it is slightly more overhead to onboard
a new one, but not much really.

This is just using http, but we can send grpc data as well in this example.

### grpc example

Start `party` api:
```
go run party/main.go --protocol grpc
```

Start the `jokeprovider`:
```
go run jokeprovider/grpc/main.go
```

Curl the API:

```
curl localhost:8090/enrol?target=jokeprovider
```

In this case we are not actually using the cloudevent sdk/system at all, but we
are using the cloudevent proto definition to be able to send generic data as
events to providers to process.

To generate the proto files:

```
cd proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative provider.proto
```

[ce]: https://cloudevents.io/
