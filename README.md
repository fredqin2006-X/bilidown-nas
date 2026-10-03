# BiliDown NAS

面向 NAS 的哔哩哔哩视频下载服务器。在浏览器中扫码登录、解析视频、提交任务，下载和合并由 NAS 后台完成。

基于 [iuroc/bilidown](https://github.com/iuroc/bilidown) 二次开发，保留上游 Git 历史与 Apache-2.0 许可证。这是独立维护的 NAS 适配项目，非上游官方发行版。

[下载镜像](https://github.com/fredqin2006-X/bilidown-nas/releases) · [部署文档](NAS部署说明.md) · [上游原始说明](docs/UPSTREAM_README.md)

![BiliDown NAS 界面预览](docs/ui-preview.jpg)

## 功能

- 扫码登录，支持普通视频、番剧/影视、合集和收藏夹解析。
- 保留上游清晰度、音频与批量下载能力；可用格式受 B 站账号权限和视频本身限制。
- 后台下载和 FFmpeg 合并，关闭浏览器后任务继续执行。
- 登录状态与任务记录持久化，下载文件直接保存在 NAS 挂载目录。
- 任务列表支持浏览器播放、下载到电脑和复制 NAS 文件路径。
- Docker 无桌面部署，普通用户运行、只读根文件系统和健康检查。
- 插画背景、半透明面板、紫色主题，适配桌面和手机。

已在绿联 DXP4800+（UGOS Pro，x86_64）验证。当前发行镜像为 **linux/amd64**，ARM NAS 尚未验证。

## 快速部署

需要 Docker 和 Docker Compose。以下步骤在 NAS 的 SSH 终端执行；按系统权限添加 `sudo`。

1. 克隆项目，并复制配置模板：

   ```sh
   git clone https://github.com/fredqin2006-X/bilidown-nas.git
   cd bilidown-nas
   cp .env.example .env
   ```

2. 编辑 `.env`：填写 NAS 的内网 IP、运行用户的 UID/GID、数据库目录和下载目录。使用 `id` 查询 UID/GID，并确保该用户可读写两个挂载目录。示例目录可改为自己 NAS 的路径。

3. 从 [Releases](https://github.com/fredqin2006-X/bilidown-nas/releases) 下载 `bilidown-nas-2.1.1-nas.3.tar.gz`，放到项目目录后导入并启动：

   ```sh
   docker load -i bilidown-nas-2.1.1-nas.3.tar.gz
   docker compose up -d --no-build
   docker compose ps
   ```

4. 浏览器访问 `http://你的NAS内网IP:8098`，使用 B 站客户端扫码登录。

也可以直接从完整源码构建，首次构建需要下载 Node/Go 和系统依赖：

```sh
docker compose build
docker compose up -d
```

## 使用与维护

下载目录固定为 `DOWNLOAD_DIR` 对应的 NAS 路径。网页设置页显示该路径；修改目录时编辑 `.env` 并重新创建容器。

所有访问者共用一个 B 站登录和任务列表，网页没有独立管理密码。设计用途为可信内网访问，请保持内网 IP 绑定。

重启保留登录和历史记录；未完成任务会标记为失败，需要重新提交。当前没有断点续传。浏览器可播放的编码由设备和浏览器决定，其他格式可以通过 NAS 文件共享使用外部播放器观看。

升级前停止下载并备份数据目录，具体见 [部署文档](NAS部署说明.md)。离线镜像不包含登录数据、数据库或下载视频。

## 开发与背景替换

前端使用 Vite、VanJS 和 Bootstrap；后端使用 Go、SQLite 和 FFmpeg。Dockerfile 包含完整构建步骤；`headless` 构建标签去除桌面托盘依赖。

替换 `client/public/background.jpg` 后重新构建即可更换背景。插画及截图中第三方内容的权利归各自权利人所有，不纳入代码的 Apache-2.0 授权，详见 [NOTICE](NOTICE)。

针对 NAS 适配的测试与完整 Docker 构建由 GitHub Actions 检查。上游部分测试依赖实时 B 站 API 或本地媒体文件，详见部署文档中的测试命令。

## 致谢与许可证

感谢 [iuroc/bilidown](https://github.com/iuroc/bilidown) 及其贡献者。NAS 适配基于提交 `6b196f57473a36ccd7c77104e5aa01ca046b7836`，修改摘要见 [NOTICE](NOTICE)。代码遵循 [Apache License 2.0](LICENSE)。
