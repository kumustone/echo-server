
# Echo Server

[中文版请点这里](README_cn.md)

Echo Server HTTP response traffic constructor for convenient testing. The echo-server provides HTTP interfaces by returning response bodies stored in files via interfaces and file downloads. It allows customization of response headers, response bodies, transfer encoding, etc. Suitable for various HTTP traffic testing scenarios, such as HTTP traffic auditing, gateway functionality testing, and frontend stubbing tests.

1. If you just want to construct response body data, you can place the response body data in the /data directory. The echo-server will return the corresponding response body based on different file suffix types.

If you simply want to simulate a response body, place the data in the /data directory, and the echo-server will return different data based on different file extensions:
1. Files with a .js suffix will return a response body with ContentType: application/javascript;charset=utf8;
2. Requests with a .json suffix will return a response body with Content-Type: application/json;charset=utf8;
3. Requests with a .html suffix will return a response body with Content-Type: text/html;charset=utf8;
4. Interfaces defined by other suffix files will return a response body with Content-Type: text/plain;charset=utf8.

2. If you want to customize the response body, place the data in the /echo directory.

For example, you can place captured response traffic in the ./echo directory and name it get_user_info.

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

You can access it via curl -v http://yourip:yourport/echo/get_user_info.

```
    |echo
    |├── get_user_info
    |└── sample_response.json
```

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

    * If you need to modify the response body in echo, you can directly modify it without rewriting the Content-Length field. The echo-server will calculate the Content-Length field based on the body.
    * If you want to add custom response headers, you can directly add the corresponding header fields in the file.
    * If you want to get the response body in chunked mode:
        curl -v http://localhost:64000/echo/get_user_info?chunk=true
    * If you want to get the response body in gzip mode:
        curl -v http://localhost:64000/echo/get_user_info?gzip=true

3. If you want to construct various file download traffic, simply place the files in the /download directory.

    You can directly open the file server through the browser at http://localhost:64000/download.

## Features

- Provides static file services from the `data`, `download`, and `echo` directories.
- Custom 404 handler displaying help information.
- Command-line options for help and version information.

## Project Structure

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

## Installation Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/echo-server.git
    cd echo-server
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Run the installation script (if any):

    ```bash
    ./install.sh
    ```

## Usage

To run the server, execute the following command:

```bash
go run main.go
```

The server will start and listen on `http://0.0.0.0:64000`.

### Command-line Options

- `-h`, `--help`: Display help information.
- `-v`, `--version`: Display version information.

### Example Requests

- Get a JSON file:

    ```bash
    curl -v http://<your-ip>:64000/data/1.json
    ```

- Get a chunked JSON file:

    ```bash
    curl -v http://<your-ip>:64000/1.json?chunk=true
    ```

- Get a gzip-encoded JSON file:

    ```bash
    curl -v http://<your-ip>:64000/data/1.json?gzip=true
    ```

## Directory Structure

- `data/`: Contains sample HTML, JavaScript, JSON, and Excel files.
- `download/`: Directory for downloadable files.
- `echo/`: Directory for files to be echoed back to the client.
- `main.go`: Main code for the Go server.
- `install.sh`: Installation script.
- `run.sh`: Script to run the server.

## License

This project uses the MIT License. For details, please see the [LICENSE](LICENSE) file.

## Contributions

Contributions are welcome! For any changes, please raise an issue or submit a pull request.

## Acknowledgments

Special thanks to all contributors and users.
