
## Private Channel Example

### 环境

+ Laravel5.8
+ Redis

### Step

首先，修改 `EventServiceProvider` 文件，在 `$listen` 添加：

```
'App\\Events\\TestPrivateBroadcastEvent' => [
    'App\\Listeners\\TestPrivateBroadcastListener',
],
```

然后执行：

```
php artisan event:generate
```

紧接着，修改 `TestPrivateBroadcastEvent` 内容如下：

```php
<?php

namespace App\Events;

use Illuminate\Broadcasting\Channel;
use Illuminate\Queue\SerializesModels;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Broadcasting\PresenceChannel;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Contracts\Broadcasting\ShouldBroadcast;

class TestPrivateBroadcastEvent implements ShouldBroadcast
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    public function broadcastOn()
    {
        return new PrivateChannel('App.User.1');
    }
}
```

> 注意这里的频道名，我们定义的是 `App.User.1` ，其真是的频道名是 `private-App.User.1` 。这个后面我们需要用得到。

之后，我们在命令行运行：

```
php artisan queue:work
```

启动队列处理进程。接下来，运行下面命令：

```
➜  laravel5.8 php artisan tinker
Psy Shell v0.9.9 (PHP 7.1.22 — cli) by Justin Hileman
>>> event(new \App\Events\TestPrivateBroadcastEvent());
=> [
     null,
   ]
>>>
```

然后我们在 redis 的服务可以看到：

```
➜  ~ docker exec -it redis1 sh
# redis-cli
127.0.0.1:6379> psubscribe *
Reading messages... (press Ctrl-C to quit)
1) "psubscribe"
2) "*"
3) (integer) 1
1) "pmessage"
2) "*"
3) "private-App.User.1"
4) "{\"event\":\"App\\\\Events\\\\TestPrivateBroadcastEvent\",\"data\":{\"socket\":null},\"socket\":null}"
```

没有问题，已经成功的进行了广播。接下来我们启动下 `go-laravel-broadcast` 服务：

```
➜  laravel-broadcasting git:(master) ✗ ./client
INFO[0000] 0.0.0.0:8890/ws
INFO[0000] redis subscribe                               channel="private-App.User.*"
```

我们还需要修改前端的内容，这里我们就在 `welcome.blade.php` 这个文件修改，内容如下：

```
<!doctype html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Laravel</title>
    </head>
    <body>

    <script type="text/javascript">
               var ws = new WebSocket("ws://127.0.0.1:8890/ws?channel=private-App.User.1");

               ws.onopen = function()
               {
                  console.log("连接成功");
               };

               ws.onmessage = function (evt)
               {
                  var received_msg = evt.data;
                  console.log(evt.data);
               };

               ws.onclose = function()
               {
                  console.log("连接已关闭...");
               };
    </script>

    </body>
</html>
```

我们在 `welcome.blade.php` 创建了一个 websocket 客户端，连接到了 `ws://127.0.0.1:8890/ws?channel=private-App.User.1` 这个地址，其中的
`channel` 的值就是我们前面 event 注册的频道名。我们先启动一个服务：

```
php artisan serve
```

然后访问 `http://127.0.0.1:8000/`，打开控制台，可以看到：

```
连接成功
(index):26 连接已关闭...
```

ws 连接被服务器主动关闭，为什么？因为 private 的 channel 需要严重用户的，所以我们需要先登录下，到 `http://127.0.0.1:8000/login` 登录下，在访问首页：

```
连接成功
3(index):21 hb
```

可以看到连接成功，没有被关闭，且受到了 `hb` 的文本消息，这个什么？这是 `go-laravel-server` 的心跳消息哦。

接下来，我们在运行下下面的命令：

```
➜  laravel5.8 php artisan tinker
Psy Shell v0.9.9 (PHP 7.1.22 — cli) by Justin Hileman
>>> event(new \App\Events\TestPrivateBroadcastEvent());
=> [
     null,
   ]
>>>
```

在浏览器的控制台上，我们可以看到：

```
连接成功
19(index):21 hb
(index):21 {"event":"App\\Events\\TestPrivateBroadcastEvent","data":{"socket":null},"socket":null}
2(index):21 hb
```

收到了 `TestPrivateBroadcastEvent` 的消息了，这样的话，我们就利用 `go-laravel-broadcast` 实现了实时通讯啦。