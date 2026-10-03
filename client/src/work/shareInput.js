/** Extract the first Bilibili URL from copied share text, keeping query parameters. */
export const normalizeShareInput = (text) => {
    const input = text.trim()
    const links = input.matchAll(/(?:https?:\/\/|\/\/)?(?:[a-z0-9-]+\.)*(?:bilibili\.com|b23\.tv)\/[^\s<>"'`“”‘’【】《》「」『』（），。！？；：、]*/gi)
    for (const match of links) {
        // Do not mistake another hostname or a URL query value for a shared link.
        if (match.index > 0 && /[\w./@=%-]/.test(input[match.index - 1])) continue
        const candidate = match[0].replace(/[)\]},.;!?]+$/, '')
        const url = new URL(candidate.startsWith('//') ? `https:${candidate}`
            : /^https?:\/\//i.test(candidate) ? candidate : `https://${candidate}`)
        if (!['bilibili.com', 'www.bilibili.com', 'm.bilibili.com', 'space.bilibili.com', 'b23.tv'].includes(url.hostname)) continue
        if (url.hostname === 'b23.tv') url.protocol = 'https:'
        if (url.hostname === 'bilibili.com' || url.hostname === 'm.bilibili.com') url.hostname = 'www.bilibili.com'
        return url.href
    }
    // Keep direct identifiers working, including the uppercase EP/SS spelling shown in the UI.
    return input.replace(/^(EP|SS)(\d+)$/i, (_, type, id) => type.toLowerCase() + id)
}
