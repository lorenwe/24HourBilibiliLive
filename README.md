# 24HourBilibiliLive

### 哔哩哔哩 直播 24小时循环推流程序

GO 语言编写的，可自行编译成各个平台的执行程序，无需前置安装任何依赖程序（ffmpeg除外）

前置条件 必须安装好 ffmpeg 切加入了 path 路径中，也就是说可以在任何目录中都可以调用到 ffmpeg

#### go 交叉编译指南

设置环境变量
```
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm64
```
#### 检查设置是否生效

go env
看到以上三项正确就说明ok了

#### 编译
```
go build main.go
```
在当前文件目录中既可以看见可执行文件 main or main.exe

将可执行文件移动到运行系统中
如 ： 
linux环境要先设置可执行权限
```
sudo chmod +x ./main
```
./main 视频文件保存目录  视频文件后缀  '推流地址'
例： ```./main /root/mp4   .mp4  'srt://...'```
注： 视频文件保存目录要从根目录开始起算，不支持相对目录，视频文件后缀如 ".mp4" 包含 "."，推流地址是srt推流地址，
取用方式后面会有解释，如需要rtmp的可以自行修改源码实现

#### srt地址说明
以B站为例，B站的srt地址有两个，服务器地址和串流密钥，其实服务器地址就已经包含了串流密钥，所以只用服务器地址就可以了
"srt://live-push.bilivideo.com:1937?streamid=#!::h=live-push.bilivideo.com,r=live-bvc/?streamname=liveid,key=密钥,schedule=srtts,pflag=1"
还可以把 包括 schedule 及以后的内容去除，则推流地址就是
"srt://live-push.bilivideo.com:1937?streamid=#!::h=live-push.bilivideo.com,r=live-bvc/?streamname=liveid,key=密钥"
因为链接中有特殊字符，所以调用推流程序传递命令行参数时需要用 英文 “单引号” 包裹传递
```
./main /root/mp4   .mp4  'srt://...'
```

此推流程序是搜索文件夹内所有符合条件的视频文件进行逐个推流，顺序需要提前通过命名排列好，全部推流完成后又会从第一个视频开始，以此达到不间断推流效果