# surge_proxy_list_gen

## 缘由

从Surge 3开始，策略组直接可以使用来自于另一个文件或者 URL 的代理声明[^surge]：

```
Group-A = fallback, policy-path=proxies.list, url = http://www.bing.com/, timeout = 2
Group-B = select, policy-path=http://example.com/proxies.txt
```

代理服务提供了一个远程列表，不过我有需求按照地区区分开来（Netflix/YouTube Premium），所以自己写个小工具暴力拆分.

## 用法

```go
go get -u -v github.com/jostyee/surge_proxy_list_gen
go run main.go -url='$remote_proxies_list' -regions='[SS] 香港,[SS] 日本,[SS] 美国,[SS] 俄罗斯' -path='$surge_icloud_drive_directory'
```

[^surge]: https://nssurge.zendesk.com/hc/zh-cn/articles/360010038714
