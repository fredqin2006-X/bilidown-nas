import assert from 'node:assert/strict'
import test from 'node:test'
import { normalizeShareInput } from '../src/work/shareInput.js'

test('extract the user share format and retain query parameters', () => {
    assert.equal(normalizeShareInput('【天星问【2025拜年纪单品】】https://www.bilibili.com/video/BV1AvfpYsE4t?vd_source=example'),
        'https://www.bilibili.com/video/BV1AvfpYsE4t?vd_source=example')
})

test('accept multiline app shares, short links and trailing punctuation', () => {
    for (const text of ['标题\nhttps://b23.tv/abc123 复制到浏览器打开', '【标题】https://b23.tv/abc123。', '(https://b23.tv/abc123)', '“https://b23.tv/abc123”']) {
        assert.equal(normalizeShareInput(text), 'https://b23.tv/abc123')
    }
})

test('normalize bare/mobile URLs without losing pages, collection or favorite IDs', () => {
    const cases = [
        ['www.bilibili.com/video/BV1AvfpYsE4t?p=2', 'https://www.bilibili.com/video/BV1AvfpYsE4t?p=2'],
        ['https://m.bilibili.com/video/BV1AvfpYsE4t', 'https://www.bilibili.com/video/BV1AvfpYsE4t'],
        ['https://bilibili.com/bangumi/play/ep123', 'https://www.bilibili.com/bangumi/play/ep123'],
        ['标题 https://space.bilibili.com/123/channel/collectiondetail?sid=456&from=share', 'https://space.bilibili.com/123/channel/collectiondetail?sid=456&from=share'],
        ['收藏 https://space.bilibili.com/123/favlist?fid=456', 'https://space.bilibili.com/123/favlist?fid=456'],
        ['http://b23.tv/ep123', 'https://b23.tv/ep123'],
        ['//www.bilibili.com/bangumi/play/ss123', 'https://www.bilibili.com/bangumi/play/ss123'],
    ]
    for (const [text, expected] of cases) assert.equal(normalizeShareInput(text), expected)
})

test('keep direct identifiers and plain URLs compatible', () => {
    for (const input of ['BV1AvfpYsE4t', 'ep123', 'ss123', 'https://www.bilibili.com/video/BV1AvfpYsE4t']) {
        assert.equal(normalizeShareInput(`  ${input}\n`), input)
    }
    assert.equal(normalizeShareInput('EP123'), 'ep123')
    assert.equal(normalizeShareInput('SS456'), 'ss456')
})

test('prefer the first actual Bilibili link over title IDs and unrelated URLs', () => {
    assert.equal(normalizeShareInput('【标题BV1other】https://example.com https://www.bilibili.com/video/BV1AvfpYsE4t https://b23.tv/other'),
        'https://www.bilibili.com/video/BV1AvfpYsE4t')
})

test('do not turn unrelated domains or nested URL parameters into Bilibili URLs', () => {
    for (const input of ['不是链接', 'https://evilbilibili.com/video/BV1AvfpYsE4t', 'https://www.bilibili.com.evil.test/video/BV1AvfpYsE4t', 'https://example.com/?url=https://www.bilibili.com/video/BV1AvfpYsE4t']) {
        assert.equal(normalizeShareInput(input), input)
    }
})
