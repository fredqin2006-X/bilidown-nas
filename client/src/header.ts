import van from 'vanjs-core'
import { now } from 'vanjs-router'
import { GLOBAL_HAS_LOGIN } from './mixin'

const { a, div, span } = van.tags

export default () => {
    const classStr = (name: string) => van.derive(() => `text-nowrap nav-link ${now.val.split('/')[0] == name ? 'active' : ''}`)

    return div({ class: 'app-header' },
        div({ class: 'app-brand' },
            span({ class: 'brand-mark', 'aria-hidden': 'true' }, 'B'),
            div({},
                div({ class: 'brand-name' }, 'Bilidown'),
                div({ class: 'brand-caption' }, '内网下载工作台'),
            ),
        ),
        div({ class: 'nav nav-underline flex-nowrap overflow-auto app-navigation', 'aria-label': '主导航' },
            div({ class: 'nav-item', hidden: () => !GLOBAL_HAS_LOGIN.val },
                a({ class: classStr('work'), href: '#/work' }, '视频解析')
            ),
            div({ class: 'nav-item', hidden: () => !GLOBAL_HAS_LOGIN.val },
                a({ class: classStr('task'), href: '#/task' }, '任务列表')
            ),
            div({ class: 'nav-item', hidden: () => !GLOBAL_HAS_LOGIN.val },
                a({ class: classStr('setting'), href: '#/setting' }, '设置中心')
            ),
            div({ class: 'nav-item', hidden: GLOBAL_HAS_LOGIN },
                a({ class: classStr('login'), href: '#/login' }, '扫码登录')
            ),
        ),
        div({ class: 'service-badge' }, span({ class: 'service-dot', 'aria-hidden': 'true' }), 'NAS · 内网访问'),
    )
}
