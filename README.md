# fake-SAUer

`fake-SAUer` 是一个简单，高效的 SAUer 打卡工具，使用 Go 语言构建。

## 如何使用？

这里我推荐一下我的使用方法，借助 `screen` 工具新建终端，当然也有其他的，供参考。

```shell
git clone https://github.com/sh1luo/fake-SAUer.git
screen -S punch
cd fake-SAUer/
go build main.go
./main
```

执行完后按组合键 `ctrl+A+D` 退出终端，按住 `ctrl` 不松手按 `A` 不松手再按 `D`。

## 版本变更

### 3.0.0-release（2021-04）

这个版本优化了整体逻辑，重构了部分代码，看起来更清晰，功能与之前无差异。

需要注意的是，现在 **不提供自动获取 UUID 的功能，需要手动在配置项中填入**。

### 2.0.0（2021-01）

这个版本使用 `json` 文件传递参数而不是原先的命令行参数。

#### Features

- 使用json文件配置信息，可读性高
- 邮件通知
- 支持多用户

### 1.0.0（2020-09）

这是第一个版本，实现了一些基本功能

#### Features

- 每日自动打卡
- 所有参数由命令行参数提供，非硬编码
- 配置简洁

## Tips

1. 如果你不确定配置项中的账号对应的密码是多少，打开这个网站尝试登录，如果可以即正确。

> https://app.sau.edu.cn/form/wap/default?formid=10

## Tutorials

[我的博客](https://kcode.icu/)