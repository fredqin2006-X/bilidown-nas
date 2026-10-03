import van from 'vanjs-core'
import { now } from 'vanjs-router'

const { div, h1, p, span } = van.tags
const copy: Record<string, [string, string]> = {
    work: ['收藏喜欢的内容', '粘贴视频链接，选择清晰度，把喜欢的视频保存到 NAS。'],
    task: ['你的下载，都在这里', '查看下载进度，播放已完成的视频，管理自己的内容。'],
    setting: ['设置你的下载空间', '管理下载目录与账号，让每一次收藏都有归处。'],
    login: ['从一次扫码开始', '登录哔哩哔哩，将喜欢的内容保存到自己的空间。'],
}

export default () => {
    const current = () => copy[now.val.split('/')[0]] ?? copy.work
    return div({ class: 'page-intro' },
        div({},
            div({ class: 'intro-eyebrow' }, span({ 'aria-hidden': 'true' }, '✦'), ' 你的私人视频空间'),
            h1(() => current()[0]),
            p(() => current()[1]),
        ),
        div({ class: 'intro-note', 'aria-hidden': 'true' }, '让喜欢的片刻', van.tags.br(), '有一个自己的家。'),
    )
}
