# 这个项目是 golang 写的。基于 web restful api 和 web 前端的音乐播放器。

项目名称 mp3-player

### web：
- 前端是 tailwind css 写的
- 使用 Alpinejs 实现局部更新

### 后端：
- 使用 gin 框架

### 数据存储
- 使用文件存储音乐库目录信息
- 每次启动遍历音乐库目录
- 每次音乐库目录变化，更新音乐库目录信息

### 功能
- 添加音乐库的目录，web 通过 api 获取路径列表进行选择
- 索引所有音乐库目录下的音乐文件
- web 获得所有音乐文件后，进行播放列表的展示
- 播放器播放音乐，暂停，上一首，下一首，随机播放，循环播放，音量控制

## 后端设计

### 用户登录
- 请求方法 POST
- 请求路径 /api/logincheck
- 请求参数
  - qiniu_access_key: 七牛云访问密钥
  - qiniu_secret_key: 七牛云秘密密钥
  - qiniu_bucket: 七牛云存储空间名称
  - qiniu_zone: 七牛云存储区域
- 相应示例
```json
正确示例
{
  "status": "success",
  "message": "登录成功",
  "data": {
    "auth_token": "xxx",
    "filelist": [
      {
        "filename": "音乐文件名",
      },
      ...
    ]
  }
}

错误示例
{
  "status": "error",
  "message": "登录失败"
}
```

### 获取音乐
- 请求方法 GET
- 请求路径 /api/stream/:songPath
- 认证 header 需要传入 auth_token
- 请求参数
  - songPath: 音乐文件路径
- 响应示例
```
c.Header("Content-Type", "audio/mpeg")
c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
c.Header("Accept-Ranges", "bytes")
```

### 上传音乐
- 请求方法 POST
- 请求路径 /api/upload
- 认证 header 需要传入 auth_token
- 请求参数
  - file: 音乐文件
- 响应示例
```json
{
  "status": "success",
  "message": "上传成功"
}

错误示例
{
  "status": "error",
  "message": "上传失败"
}
```

### 删除音乐
- 请求方法 DELETE
- 请求路径 /api/delete/:songPath
- 认证 header 需要传入 auth_token
- 请求参数
  - songPath: 音乐文件路径
- 响应示例
```json
{
  "status": "success",
  "message": "删除成功"
}

错误示例
{
  "status": "error",
  "message": "删除失败"
}
```


## 测试
- 运行测试
```
go test ./tests
go test -v -count=1 ./tests/ -run TestCheckQiniuAccess
```
