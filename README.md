# wechatbot

> 本项目是 fork 他人的项目来进行学习和使用，请勿商用，可以下载下来做自定义的功能
> 最近ChatGPT异常火爆，本项目可以将个人微信化身GPT机器人，
> 项目基于[openwechat](https://github.com/eatmoreapple/openwechat) 开发。

## 目前实现了以下功能

* 默认使用`gpt-3.5-turbo`模型，可选配置`text-davinci-003`，gpt4还在捣鼓中，相信很快就会到来~
* 支持配置base_url和proxy，方便国内用户使用。二者配置其一即可，参数解释详见[配置文件说明](#配置文件说明)
* 支持缓存对话上下文，可配置缓存时间，默认缓存1分钟
* 支持通过指令清空上下文，可配置清空指令
* 支持群聊@回复
* 支持私聊回复
* 支持私聊回复前缀设置
* 支持配置好友添加是否自动通过

## 效果展示
![](https://image.aayu.today/uploads/2023/04/03/202304030328786.png)

# 实现机制
目前机器人有两种实现方式
* 逆向功能，扒取官网API，通过抓取cookie获取GPT响应信息，`优点：`效果与官网一致，`缺点：`cookie会过期需要不定时更新。
* 基于openai官网提供的API，`优点`：模型以及各种参数可以自由配置，`缺点：`API收费，新账号有5美元免费额度，大约可以支持200万个英文单词的对话，或125万个中文单词的对话。

> 本项目基于第二种方式实现，模型之间具体差异可以参考[官方文档](https://beta.openai.com/docs/models/overview), 详细[参数示例](https://beta.openai.com/examples) 。

# 常见问题
* 如无法登录 login error: write storage.json: bad file descriptor 删除掉storage.json文件重新登录。
* 如无法登录 login error: wechat network error: Get "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage": 301 response missing Location header 一般是微信登录权限问题，先确保PC端能否正常登录。
* 其他无法登录问题，依然尝试删除掉storage.json文件，结束进程(linux一般是kill -9 进程id)之后重启程序，重新扫码登录，(如为docket部署，Supervisord进程管理工具会自动重启程序)。
* 机器人无法正常回复，检查ApiKey能否正常使用，控制台日志中有详细错误信息
* linux中二维码无法扫描，缩小命令行功能，让二维码像素尽可能清晰。（无法从代码层面解决）
* 机器人一直答非所问，可能因为上下文累积过多。切换不同问题时，可以先发送配置的清空指令`session_clear_token`，清空上下文，即可重新开始对话~

# 使用前提

* 有openai账号，并且创建好api_key。
* 微信必须实名认证。

# 注意事项

* 项目仅供娱乐，滥用可能有微信封禁的风险，请勿用于商业用途。
* 请注意收发敏感信息，本项目不做信息过滤。

# 使用docker运行

你可以使用docker快速运行本项目。

`第一种：基于环境变量运行`

```sh
# 运行项目，环境变量参考下方配置说明
$ docker run -itd --name wechatbot --restart=always \
 -e APIKEY=换成你的key \
 -e AUTO_PASS=false \
 -e SESSION_TIMEOUT=60s \
 -e MODEL=text-davinci-003 \
 -e MAX_TOKENS=512 \
 -e TEMPREATURE=0.9 \
 -e REPLY_PREFIX=我是来自机器人回复: \
 -e SESSION_CLEAR_TOKEN=清空会话 \
 ylsislove/wechatbot:latest

# 查看二维码
$ docker exec -it wechatbot bash 
$ tail -f -n 50 /app/run.log 
```

`第二种：基于配置文件挂载运行`

```sh
# 复制配置文件，根据自己实际情况，调整配置里的内容
$ cp config.dev.json config.json  # 其中 config.dev.json 从项目的根目录获取

# 运行项目
$ docker run -itd --name wechatbot -v `pwd`/config.json:/app/config.json ylsislove/wechatbot:latest

# 查看二维码
$ docker exec -it wechatbot bash 
$ tail -f -n 50 /app/run.log 
```

其中配置文件参考下边的配置文件说明。

# 快速开始

`第一种：直接下载二进制(适合对编程不了解的同学)`

> 非技术人员请直接下载release中的[压缩包](https://github.com/869413421/wechatbot/releases) ，请根据自己系统以及架构选择合适的压缩包，下载之后直接解压运行。

下载之后，在本地解压，即可看到可执行程序，与配置文件：

```bash
# windows
1.下载压缩包解压
2.复制文件中config.dev.json更改为config.json
3.将config.json中的api_key替换为自己的
4.双击exe，扫码登录

# linux
$ tar xf wechatbot-v0.0.2-darwin-arm64.tar.gz
$ cd wechatbot-v0.0.2-darwin-arm64
$ cp config.dev.json # 根据情况调整配置文件内容
$ ./wechatbot  # 直接运行

# 如果要守护在后台运行
$ nohup ./wechatbot &> run.log &
$ tail -f run.log
```

`第二种：基于源码运行(适合了解go语言编程的同学)`

````
# 获取项目
$ git clone https://github.com/ylsislove/wechatbot.git

# 进入项目目录
$ cd wechatbot

# 复制配置文件
$ copy config.dev.json config.json

# 启动项目
$ go run main.go
````

# 配置文件说明

```json
{
  "api_key": "your api key",
  "auto_pass": true,
  "session_timeout": 60,
  "max_tokens": 1024,
  "model": "gpt-3.5-turbo",
  "temperature": 1,
  "reply_prefix": "来自机器人回复：",
  "session_clear_token": "清空会话",
  "base_url": "https://api.openai.com/v1/",
  "request_timeout": 60,
  "proxy": ""
}
```

* api_key：openai api_key
* auto_pass：是否自动通过好友添加
* session_timeout：会话超时时间，默认60秒，单位秒，在会话时间内所有发送给机器人的信息会作为上下文。
* max_tokens: GPT响应字符数，最大2048，默认值1024。max_tokens会影响接口响应速度，字符越多响应越慢。
* model: GPT选用模型，默认`gpt-3.5-turbo`，具体选项参考官网训练场
* temperature: GPT热度，0到1，默认1。数字越大创造力越强，但更偏离训练事实，越低越接近训练事实
* reply_prefix: 私聊的回复前缀
* session_clear_token: 情况上下文的指令，默认`清空会话`
* base_url：API请求地址，默认官网API地址，国内用户可设置为自己搭建的代理地址
* request_timeout：请求超时时间，默认60秒
* proxy: 代理地址，如果需要使用代理，可以在此配置代理地址。例如：`http://127.0.0.1:18081`
