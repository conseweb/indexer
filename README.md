# indexer

`go run cli/main.go`

### 获取设备列表
Get 	"/indexer/devices"

### 获取设备上的所有文件
Get 	"/indexer/devices/:device_id"

### 使设备上线可用，body如下json，下同
Post 	"/indexer/devices/:device_id/online"

### 使设备下线
Post 	"/indexer/devices/:device_id/offline"

### 更新设备上的文件索引
Post 	"/indexer/devices/:device_id/files"

### 更新文件索引信息
Post 	"/indexer/files"

### 获取单个文件索引信息
Get 	"/indexer/files/:file_id"

### 删除文件索引信息
Delete 	"/indexer/files/:file_id"


```
[
	{
		"id":"aaaaaaaa",
		"address":"127.0.0.1:1234"
	},

	{
		"id":0,
		"device_id":"aaaaaaaa",
		"path":"/a",
		“size":1234,
		"has":"fffff",
		"created":"0001-01-01T00:00:00Z",
		"updated":"0001-01-01T00:00:00Z"
	}
]
```