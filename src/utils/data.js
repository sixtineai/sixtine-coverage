import config from '@/config'

let codeFilesMapCache = null

function isGzBase64(str) {
    if (!str || str.length < 4) return false
    try {
        const raw = atob(str.slice(0, 4))
        return raw.charCodeAt(0) === 0x1f && raw.charCodeAt(1) === 0x8b
    } catch {
        return false
    }
}

async function decompressGzBase64(b64) {
    const bin = atob(b64)
    const bytes = new Uint8Array(bin.length)
    for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i)
    const ds = new DecompressionStream('gzip')
    const writer = ds.writable.getWriter()
    writer.write(bytes)
    writer.close()
    const reader = ds.readable.getReader()
    const chunks = []
    while (true) {
        const { done, value } = await reader.read()
        if (done) break
        chunks.push(value)
    }
    const total = chunks.reduce((s, c) => s + c.length, 0)
    const result = new Uint8Array(total)
    let offset = 0
    for (const c of chunks) { result.set(c, offset); offset += c.length }
    return new TextDecoder().decode(result)
}

export async function resolveStringOrGz(value) {
    if (!value) return null
    if (isGzBase64(value)) return decompressGzBase64(value)
    return value
}

export async function getCodeFilesMap() {
    if (codeFilesMapCache !== null) return codeFilesMapCache

    const raw = config.codeFilesMap
    if (!raw) {
        codeFilesMapCache = null
        return null
    }

    try {
        const str = await resolveStringOrGz(raw)
        codeFilesMapCache = JSON.parse(str)
        return codeFilesMapCache
    } catch (e) {
        console.error('Failed to parse code files map:', e)
        codeFilesMapCache = null
        return null
    }
}

export function hasInlineCodeFilesMap() {
    return !!config.codeFilesMap
}

export async function resolveInlineCoverage() {
    const raw = config.coverageData
    if (!raw) return null
    return resolveStringOrGz(raw)
}

export async function resolveInlineRouteCoverage() {
    const raw = config.routeCoverageData
    if (!raw) return null
    try {
        return await resolveStringOrGz(raw)
    } catch (e) {
        console.error('Failed to resolve inline route coverage:', e)
        return null
    }
}
