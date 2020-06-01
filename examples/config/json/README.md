# JSON文件解析

Go 标准库提供了 encoding/json 包来解析 json 文件。

以例子来说明 encoding/json 包的用法。

配置文件如下：

```
{
	"port": 10666,
	"mysql": {
		"url": "(27.0.0.1:3306)/biezhi",
		"username": "root",
		"password": "123456"
	}
}
```
代码结构如下：
```
type MysqlConfig struct {
	URL string
	UserName string
	Password string
}

type Config struct {
	PorT int
	MySQl MysqlConfig
}
```

可以看出：

1. 代码结构跟 json 文件中的字段并不一致，但是在不区分大小写的情况下，字母还是一致的。
2. struct 里面的字段需要首字母大写，可以被导出。
3. struct 里面的字段也可以跟 json 文件中的不一致，需要加上 tag，如下所示，也可以解析出来，具体的示例在 main2.go 中。



```
type MysqlConfig struct {
	u string            `json:"url"`
	un string 	        `json:"username"`
	pw string           `json:"password"`
}
```



