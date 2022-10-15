# Simple Blog Community

​		本项目为参考 CSDN 搭建的一个简易博客社区，包括登录注册、关键字查询、分类筛选、文章的增删查改以及收藏与关注功能。请按照以下步骤，让项目在跑起来吧！

## 1. 环境要求

在开始之前，请保证您安装好以下环境（本项目所使用的版本仅供参考）：

|         | 版本查看语句      | 本项目使用的版本 |
| ------- | ----------------- | ---------------- |
| Go      | `go version`      | 1.17.6           |
| Node.js | `node --version`  | 16.14.0          |
| Vue.js  | `vue --version`   | 2.9.6            |
| MySQL   | `mysql --version` | 8.0.28           |

## 2. 数据库准备

在mysql上新建数据库blog，新建分类表categories，字段名为id和category_name，并插入数据（其余表自动生成，无需创建）：

![image1](https://s3.bmp.ovh/imgs/2022/10/15/f538c8d486355413.png)

打开blog_server/common/database.go，在InitDB函数中修改mysql用户名及密码：

![image2](https://s3.bmp.ovh/imgs/2022/10/15/c45c52e59b1d9322.png)

## 3. 启动项目

从终端进入blog_server，输入以下语句启动后端：

```
go run main.go
```

从前端进入blog_client，输入以下语句启动前端：

```
npm run dev
```

[项目展示视频链接](https://www.bilibili.com/video/BV16e4y1n75z/?spm_id_from=333.999.0.0&vd_source=45ccd2342317a2b6f76d814e6cb88aca)

