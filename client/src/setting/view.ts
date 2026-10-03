import van from 'vanjs-core'
import { SettingRoute } from '.'
import { saveFields } from './data'
import { SERVER_INFO } from '../mixin'

const { a, button, div, input } = van.tags

export const SaveFolderSetting = (route: SettingRoute) => {
    const saveFolder = route.fields.download_folder
    const folderPickerDisabled = van.state(false)
    const buttonText = '保存'

    return div({ class: 'input-group' },
        div({ class: 'input-group-text' }, '下载目录'),
        input({
            class: 'form-control',
            value: () => SERVER_INFO.val.serverMode ? SERVER_INFO.val.downloadHostPath : saveFolder.val,
            readOnly: () => SERVER_INFO.val.serverMode,
            oninput: event => saveFolder.val = event.target.value,
        }),
        button({
            class: 'btn btn-success', hidden: () => SERVER_INFO.val.serverMode, onclick() {
                folderPickerDisabled.val = true
                saveFields([
                    ['download_folder', saveFolder.val]
                ]).then(message => {
                    alert(message)
                }).catch(error => {
                    if (error instanceof Error) alert(error.message)
                }).finally(() => {
                    folderPickerDisabled.val = false
                })
            }, disabled: folderPickerDisabled
        }, buttonText)
    )
}