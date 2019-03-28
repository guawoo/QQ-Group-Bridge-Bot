## QQ群桥接机器人
----
本机器人使用酷Q平台的[http插件](https://cqhttp.cc/docs/4.8/#/)实现

主要功能有：

1. 广播消息
2. 群内屏蔽QQ号，不显示其广播的消息。
3. 群内屏蔽Q群号，不显示其群内成员广播的消息。
4. 自由加入和退出广场。

## 框架使用方法：

框架主骨架为observer包和broadcast包。

observer包中的observer是一个http服务器，用于侦测发来的QQ消息。
broadcast包中的station是一个多线程处理器，用于发送QQ消息和处理需要立即返回的任务。

broadcast包中的action封装了常用的酷Q http的SDK。

## 消息处理方法

使用observer中的RegiestMsgHandle注册要处理的消息的相关函数。
它的三个参数分别是：
* pattern 匹配消息用的正则表达式
* msgtype 消息的类型（自己定义，也可以使用代码中默认的）
* hanlde 消息处理函数

## 已实现的功能使用
（注意：先登录酷Q，再运行本程序）

* 在群内输入#join 加入广场，广场内的群可以收到广播。
* 在群内输入#leave 退出广场。
* 在群内输入#list 显示已加入广场的群。
* 在群内输入<<接要广播的文字 即可在加入广场中的群内广播消息。
