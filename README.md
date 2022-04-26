# go-kit

This is my way of seeing go-kit to implement a simple health. I'm using go-kit to "address incidental complexity".

Solving:
- circuit breaking
- rate limit
- instrumentation 
- logging

Honorable mention (can be added as well):
- tracing (provider not included)
- monitoring (provider not included)

Run it by executing (i like port 4041, use whatever you like):
```
go build .
# care, this is running as a daemon, if you execute it twice * boom *
HTTP_PORT=4041 ./go-kit-example & 
curl http://localhost:4041/healthz
```

File hierarchy can be address like: 

Great if you have a really small service
endpoints
    - profiles
    - users
    - images
services
    - profiles
    - users
    - images
transports 
    - http_profiles
    - grpc_profiles
    - http_users
    - http_images

If you have many resources on one service (don't recommend though):
profiles
    - service
    - endpoints
    - transport(s)
users
    - service
    - endpoints
    - transport(s)
images
    - service
    - endpoints
    - transport(s)

