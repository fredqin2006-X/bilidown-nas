/// <reference types="vite/client" />
import van from 'vanjs-core'
import Header from './header'
import Hero from './hero'
import Work from './work'
import Task from './task'
import Login from './login'
import Setting from './setting'
import _Error from './error'
import { redirect } from 'vanjs-router'
import { GLOBAL_HIDE_PAGE } from './mixin'
import 'bootstrap/dist/css/bootstrap.min.css'
import './scss/index.scss'
import { PlayerModalComp } from './task/playerModal'

const { div } = van.tags

redirect('home', 'work')

van.add(document.body,
    div({ class: 'container app-shell vstack', hidden: GLOBAL_HIDE_PAGE },
        Header(),
        Hero(),
        Work(),
        Task(),
        Login(),
        Setting(),
        div({ class: 'app-footer' }, 'Bilidown · 视频保存在 NAS，浏览器随时访问'),
    ),
    _Error()
)
