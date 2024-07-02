
# Echo Server

[Click here for the English version](README.md)


Echo Server HTTP响应流量流量构造器，方便测试使用。echo-server会讲用户存储在文件中的响应体以接口，文件下载等方式提供http接口，可以方便的自定义接口的响应头，响应体，传输编码等。适用需要于各种 HTTP 流量测试场景，如 HTTP 流量审计、网关功能测试和前端打桩测试等。


1. 如果你只是想构造响应体数据，可以把响应体数据放到/data目录下文件下。echo-server会根据不同的文件后缀类型返回对应的影响体。

如果你只是想模拟一个简单的响应体数据，把数据放到data目录下，echo-server根据不同文件名返回不同的数据:
1. .js 后缀的文件，以ContentType: application/javascript;charset=utf8 返回响应体;
2. .json 后缀的请求, 以Content-Type:application/json;charset=utf8 返回响应体;
3. .html 后缀的请求，以Content-Type:text/html;charset=utf返回响应体;
4. 其它后缀文件定义的接口， 以Content-Type:text/plain;charset=utf返回响应体;

2. 如果你想自定义响应体同时，把数据放到/echo目录下面。

以echo中的文件为例，你可以把抓取到的响应流量放到./echo目录下，并名为get_user_info.

```
HTTP/1.1 200 OK
Content-Type: application/json;charset=utf8
Date: Wed, 20 Mar 2024 01:58:06 GMT
Content-Length: 1618
 

{
"status": 1,
"data": {
"phone_number": "010-12345678",
"bus_number": "京A12345",
"org_code": "12345678-X"
}
}
```

你可以通过curl -v http://yourip:yourport/echo/get_user_info来访问它。

```
curl -v http://localhost:64000/echo/get_user_info
*   Trying 172.16.43.236:64000...
* Connected to 172.16.43.236 (172.16.43.236) port 64000
> GET /echo/login HTTP/1.1
> Host: 172.16.43.236:64000
> User-Agent: curl/8.6.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Sat, 22 Jun 2024 01:11:26 GMT
< Content-Length: 113
< Content-Type: text/plain; charset=utf-8
<
{
"status": 1,
"data": {
"phone_number": "010-12345678",
"bus_number": "京A12345",
"org_code": "12345678-X"
}
}
```

Note

    * 如果因为某些需要你想修改echo中数据响应体，可以直接修改，无需重写Content-Length字段。echo-server会根据body来计算Content-Length字段。

    * 如果你想增加自定义的响应头。可以直接在文件中添加相应头字段。

    * 如果你想以chunk的方式得到响应体:
        curl -v http://localhost:64000/echo/get_user_info?chunk=true

    * 如果你想以gzip的方式得到响应体：
        curl -v http://localhost:64000/echo/get_user_info?gzip=true


3. 如果你想构造各种文件下载的流量，直接把文件放到/download目录下。
    
    通过浏览器可以直接打开文件服务器.  http://localhost:64000/download.



## 功能特性

- 提供来自 `data`、`download` 和 `echo` 目录的静态文件服务。
- 自定义 404 处理程序，显示帮助信息。
- 命令行参数，提供帮助和版本信息。

## 项目结构

```
.
├── data
│   ├── 1.html
│   ├── 1.js
│   ├── 1.json
│   └── 1.xlsx
├── download
│   └── 1.txt
├── echo
│   └── 1.json
├── echo-server
├── echo-server.service
├── go.mod
├── go.sum
├── install.sh
├── main.go
├── README.md
└── run.sh
```

## 安装步骤

1. 克隆仓库：

    ```bash
    git clone https://github.com/yourusername/echo-server.git
    cd echo-server
    ```

2. 安装依赖：

    ```bash
    go mod tidy
    ```

3. 运行安装脚本（如果有）：

    ```bash
    ./install.sh
    ```

## 使用方法

要运行服务器，执行以下命令：

```bash
go run main.go
```

服务器将启动并监听 `http://0.0.0.0:64000`。

### 命令行参数

- `-h`, `--help`：显示帮助信息。
- `-v`, `--version`：显示版本信息。

### 示例请求

- 获取 JSON 文件：

    ```bash
    curl -v http://<your-ip>:64000/data/1.json
    ```

- 获取分块传输的 JSON 文件：

    ```bash
    curl -v http://<your-ip>:64000/1.json?chunk=true
    ```

- 获取 gzip 编码的 JSON 文件：

    ```bash
    curl -v  http://<your-ip>:64000/data/1.json?gzip=true
    ```

## 目录结构

- `data/`：包含示例 HTML、JavaScript、JSON 和 Excel 文件。
- `download/`：可下载文件的目录。
- `echo/`：要回显给客户端的文件目录。
- `main.go`：Go 服务器的主代码。
- `install.sh`：安装脚本, 会
- `run.sh`：运行服务器的脚本。

## 许可证

此项目使用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 贡献

欢迎贡献！如有任何更改，请提出 issue 或提交 pull request。

## 鸣谢

特别感谢所有贡献者和用户。
