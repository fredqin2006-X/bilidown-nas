# BiliDown NAS 部署说明

镜像版本：`bilidown-nas:2.1.1-nas.4`，架构 `linux/amd64`。上游基线：iuroc/bilidown 提交 `6b196f57473a36ccd7c77104e5aa01ca046b7836`。

## 配置

复制 `.env.example` 为 `.env`，按实际 NAS 修改：

| 变量 | 用途 |
| --- | --- |
| `LAN_IP` | NAS 内网 IPv4 地址，必须填写；HTTP 端口只绑定此地址 |
| `WEB_PORT` | 网页端口，默认 8098 |
| `PUID` / `PGID` | 运行用户的数字 UID/GID，用 `id` 查询 |
| `DATA_DIR` | NAS 数据库目录，挂载到容器 `/data` |
| `DOWNLOAD_DIR` | NAS 下载目录，挂载到容器 `/downloads` |

示例目录为 `/volume3/docker/BiliDown/data` 和 `/volume2/Download`，可按自己的 NAS 修改。部署前创建这两个目录，并赋予运行用户读写权限。

数据库为数据目录中的 `data.db`，保存 B 站登录和任务信息。不要提交 `.env`、数据库或下载视频到 Git 仓库。

所有访问者共用一个 B 站登录和任务列表，网页没有独立管理密码，供可信内网使用。无需公网域名、端口转发或反向代理。

## 导入发行镜像

从 [Releases](https://github.com/fredqin2006-X/bilidown-nas/releases) 下载镜像备份，放到项目目录。在 NAS SSH 中执行，必要时添加 `sudo`：

```sh
docker load -i bilidown-nas-2.1.1-nas.4.tar.gz
docker compose up -d --no-build
docker compose ps
docker compose logs --tail=100 bilidown
```

浏览器访问 `http://LAN_IP:WEB_PORT`（替换为配置值）。镜像备份不包含数据库和视频；恢复这些内容需要单独备份挂载目录。

## 从源码构建

Git 仓库提供完整源码构建，包含网页、Go 无界面服务及 FFmpeg：

```sh
docker compose build
docker compose up -d
```

Compose 默认使用清华 Debian 镜像加速并保留 Debian 签名校验。此设置仅作用于容器构建，不修改 NAS 系统软件源。需要原始 Debian 源时修改 compose.yaml 的构建参数。参考：[清华 Debian 镜像帮助](https://mirrors.tuna.tsinghua.edu.cn/help/debian/)。

发行附件 `bilidown-nas-build.tar.gz` 另带已编译的 Linux 二进制和静态资源，可解压后使用预构建 Dockerfile：

```sh
docker build --network=host \
  --build-arg DEBIAN_MIRROR=http://mirrors.tuna.tsinghua.edu.cn/debian \
  --build-arg DEBIAN_SECURITY_MIRROR=http://mirrors.tuna.tsinghua.edu.cn/debian-security \
  -f Dockerfile.prebuilt -t bilidown-nas:2.1.1-nas.4 .
docker compose up -d --no-build
```

`Dockerfile.prebuilt` 依赖构建附件中未纳入 Git 的 `build/` 和 `server/static/`；普通 Git 克隆请使用完整 `Dockerfile`。

## 使用与升级

扫码登录后提交下载任务，文件保存到 `DOWNLOAD_DIR`。关闭浏览器不会停止任务。任务列表可播放浏览器支持的媒体，也可下载到当前电脑；其他编码可通过 NAS 文件共享和外部播放器播放。

文件夹按钮显示可复制的 NAS 文件路径。设置页下载目录固定为映射目录；改变目录需修改 `.env` 并重新创建容器。NAS IP 改变时同步修改 `LAN_IP`。

升级前先完成或停止当前下载、停止容器，再备份 `DATA_DIR` 并记录旧镜像标签。视频按需单独备份：

```sh
docker compose stop
# 此时备份 DATA_DIR；随后导入新镜像并更新 compose.yaml 标签
docker compose up -d --no-build
```

停止服务可用 `docker compose stop` 或 NAS Docker 应用。重启保留登录与历史记录，但未完成下载标记为失败，需要重新提交，当前没有断点续传。

容器普通用户运行、根文件系统只读，不需要特权模式、核显设备或 Docker socket。只挂载数据和下载目录。首次启动失败时查看容器日志，并检查目录权限、端口占用及内网地址是否正确。

## 适配修改

- `headless` 构建标签去除桌面托盘依赖，缺少 FFmpeg 时启动报错退出。
- 数据库、静态文件、监听地址和下载根目录可通过环境变量配置。
- 文件下载只接受已完成任务 ID，限制到下载根目录并检查符号链接；支持 HTTP Range。
- 服务模式禁用退出软件、打开桌面文件管理器和更改下载根目录。
- `/healthz` 健康检查、停止信号、日志轮转与自动重启。
- 下载响应及时关闭，错误状态不写入媒体文件；进度读取使用锁保护的快照，未知文件长度不会生成非法 JSON。
- 镜像静态目录权限为 755、文件为 644，避免 NAS 解压产生 700 权限后网页资源返回 403。
- nas.3 增加插画背景、半透明面板、紫色主题和手机布局；保留原操作事件与后端功能。

## 本地开发与测试

前端使用 Node.js 22、pnpm 9.15.4，后端使用 Go 1.23 和 FFmpeg：

```sh
cd client
pnpm install --frozen-lockfile
pnpm build
cd ../server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags headless -trimpath -ldflags='-s -w' -o ../build/bilidown-linux-amd64 .
```

上游部分测试依赖本地媒体文件或实时 B 站 API。NAS 适配测试可独立运行：

```sh
cd server
go test -tags headless ./... -run 'TestHeadlessDatabaseStartup|TestResolveDownloadFile|TestServer|TestTaskDownload'
go vet -tags headless ./...
```

nas.4 增加分享文字提取：直接粘贴 B 站复制的“标题 + 链接”即可。支持 b23 短链、手机端链接、换行和末尾标点，保留链接参数及原有 BV/EP/SS 输入。前端测试：在 `client` 目录执行 `pnpm test`。

## 已验证范围（2026-10-03）

- 网页构建、Go 无界面 Linux 编译、适配单元测试及 go vet。
- NAS Linux 符号链接越界限制、固定下载目录与 HTTP Range。
- 独立测试容器使用生成的短音视频，验证后台下载、FFmpeg 合并、中文文件名、完整文件和 Range 读取。
- 重启后的数据库、任务记录与中断任务状态。
- 正式容器健康检查、内网页面、JavaScript/CSS 资源及背景图片。
- 实际扫码登录、真实视频下载功能经部署用户验收。
- 桌面及 390px 手机解析页、任务列表与设置页布局，长文件名未横向溢出。

其他 NAS 型号、ARM 架构和公网部署尚未验证。
