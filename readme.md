## web-grpc-video-chat

Hello,

This is my personal playground where I'd implemented video chat using web grpc and some service side logic.

The happy path will be like described below.

Happy path:
```
(UserA) => (Browser) => (Video/Audio Capture) => (Go backend service) => (Browser) => (UserB)
```

#### PS: We use Envoy as tls terminator, router and WebGRPC 2 gRPC transport.

### C4 Model
![](mac-video-chat.jpg)
