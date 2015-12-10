# HttpMonitor
监控Web站点，异常情况将发送邮件告知管理员


## 如何使用？
一.下载并编译  
1.[下载](https://github.com/alienxt/http-keep-alive/archive/master.zip)  
2.解压  
3.go build  

二.配置config.json  
1.在urls里配置你要监控的站点，url为要监控的站点url，interval为监控间隔，timeout为响应超时的时间，title为告警的标题  
2.设置你的邮箱服务器及用户名密码  
3.设置接收告警的邮件  
