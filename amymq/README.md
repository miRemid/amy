# config.json配置
config.json文件放在同目录文件的config目录中
```json
{
    // 配置消息管道缓冲大小
    "amymq":{
        "channel":{
            // 主管道
            "main": 50,
            "message": 20,
            "notice": 10,
            "request": 10
        }
    },
    // 配置命令开头格式，不符合格式将过滤
    "cmds":["!", "！", "#"],
    // CQHTTP插件地址端口
    "cqhttp_url": "127.0.0.1",
    "cqhttp_port": 5700,
    // 上报地址模板
    "post_url": "http://127.0.0.1:8080/api/coolq",
    "message":{
        // 上报方法(目前强制post)
        "method": "post",
        // 上报地址，默认为post_url
        "url": "",
        // 上报路由，默认为message
        "router": "/message"
    },
    "notice":{
        "method": "post",
        "url": "http://api.qqbot.hahaha",
        "router": "/notice"    
    },
    "request":{
        "method": "post",
        "url": "",
        "router": "/request"    
    },
    // 上报超时时间
    "timeout": 10000,
    // 设置日志文件夹，会自动生成amy.log日志文件，默认当前路径的log文件夹
    "log_path": "./log",
    // CQHTTP配置项
    "secret": "amy"
}
```