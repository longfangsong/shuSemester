# API Reference

## 模型

### 学期

```json
{
   id: 	     学期id
   start:    学期开始日期，日期时间类型
   end:      学期结束日期，日期时间类型
   name:     学期名称，字符串
   holidays: 假期列表
}
```

### 假期

```json
{
   id: 	   假期id
   start:  假期开始日期，日期时间类型
   end:    假期结束日期，日期时间类型
   name:   假期名称，字符串
   shifts: 调休列表
}
```

### 调休

```json
{
   id:        调休id
   rest_date: 调休补休日期，日期时间类型
   work_date: 调休补工作日期，日期时间类型
}
```

## cli

cli用于生成管理员token，

使用方法如下：

```shell
>>> ./cli [你的学生证号]
<<< [一个JWT Token]
```

## web api

- `GET /ping`

  检查服务是否可用，应该直接返回`pong`。

- `GET /semester?date=[一个日期，格式为js标准格式]`

  返回日期对应的学期。

- `GET /semester?date=now`

  返回今天所在的学期。

- `POST /semester`

  添加或修改学期。

  - `body`

    json格式的学期信息。
    
    例如
    ```json
    {
      "name":"2018-2019夏季学期",
      "start":"2019-07-16T16:00:00.000Z",
      "end":"2019-09-01T16:00:00.000Z",
      "holidays":[]
    }
    ```

  - `head`

    - `Authorization`

      `Bearer [你的JWT Token]`。

      