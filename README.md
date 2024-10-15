# 快速开始
```
git clone https://github.com/LingCode97/ssh-chatroom.git
cd ssh-chatroom
go build
./main
```

# 功能说明
运行项目后，通过``ssh username@127.0.0.1 -p 2222``加入聊天室。除基础聊天外，还支持以下指令：
* #user:查看当前房间在线用户
* #leave:离开当前聊天室
* #rooms:查看所有聊天室
* #enter room_name:进入指定聊天室

其他功能持续更新...