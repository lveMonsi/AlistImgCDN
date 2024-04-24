# AlistImgCDN
## 点击一键部署
<p><a href="https://vercel.com/new/clone?repository-url=https://github.com/lveMonsi/AlistImgCDN" target="_blank" rel="noopener noreferrer"><img loading="lazy" src="https://vercel.com/button" alt="Deploy with Vercel"></a></p>

## 部署步骤
1. 点击上面的一键部署按钮，部署到 Vercel 项目中
2. 配置环境变量，添加一个环境变量，Key 为 `url` ，Value 为你的 Alist 域名
3. 点击部署，等待云函数部署完成
4. 为该项目配置一个域名

## 使用教程
1. 复制你想要静态化的 Alist 文件下载链接
```
https://你的Alist域名/d/root/nyancat.gif
```
2. 取其主域名后的路径参数，拼接到 `该项目域名 + /img/` 后，如下
```
https://项目的域名/img/d/root/nyancat.gif
```
3. 直接使用该链接即可用作图片直链，访问可发现其为静态资源，非 Alist 提供的下载直链。

## 项目原理
使用 Gin 框架，在访问资源时解析路径参数，查询缓存 map 中是否已缓存该资源，若已缓存则判断是否过期，否则调用 golang.http 库直接回源获取对应图片资源并存入缓存 map ，之后在缓存冷却期间访问该资源将直接返回缓存 map 中的图片数据，实现白嫖 Vercel 全球 CDN。

## 弊端
1. 每个资源第一次访问的速度都约等于直接访问 Alist 的直链的速度，相当于 Vercel 只优化了高峰访问期对 Alist 端的流量缓解。
2. 实际部署起来才发现，Vercel 云函数有空闲内存回收机制，在长时间不访问该项目后，项目资源会被回收，导致缓存失效重新去 Alist 获取资源。若需要足够时间的缓存，需要借助云函数以外的数据储存，避免缓存丢失。