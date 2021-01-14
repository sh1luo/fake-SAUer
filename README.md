# fake-SAUer

`fake-SAUer` is a Simple configuration, efficient and reliable punching tool.

## Changelog

### 2.0.0（2021-01）

This version has been changed to configure information using json other than command line parameters.

#### Features

- Configure by json file.
- Email results notification.
- **Support multiple users.**

### 1.0.0（2020-09）

#### Features

- Automatically punch every day.
- Configure before running,not hard-coded in the program.
- Simple configuration.

#### Use

```shell
go run main.go signup -a=Sno -w=passwd -p=phone_number -n=your_name -r=province -c=city -o=college
```

#### Tips

1. If you are not sure about the password,you can try to log in to the website below,if okay,it means correct.

> https://app.sau.edu.cn/form/wap/default?formid=10

## Tutorials

[My Blog](https://kcode.icu/)



