## web-grpc-video-chat

Hello,

This is my personal playground where I want to implement video chat using grpc and server side logic in Go.

My happy path will be like described further.

Happy path:
```
(UserA) => (Browser) => (Video/Audio Capture) => (Go backend service) => (Browser) => (UserB)
```

First what we do is a c4 model, in our case c3, it's very helpful and save a lot of time.

![](mac-video-chat.jpg)


After joining and when users can see and hear each other, the PoC ends.

Main goals is to find out resources needed and opencv operations.

I do this strictly for self-education and a pleasant pastime. I do not give any guarantees,
and you are free to use this code if you suddenly need it. ;)