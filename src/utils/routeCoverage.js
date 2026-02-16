import config from '@/config'

export function parseRouteCoverage(text) {
    const lines = text.split('\n')
    let currentTag = 'Unknown'

    const routeList = []
    const serviceMap = {}

    for (const raw of lines) {
        const line = raw.trim()
        if (!line || line.startsWith('mode:')) continue

        if (line.startsWith('#')) {
            currentTag = line.slice(1).trim()
            continue
        }

        const firstSpace = line.indexOf(' ')
        if (firstSpace < 0) continue
        const method = line.slice(0, firstSpace)

        const rest = line.slice(firstSpace + 1)
        const secondSpace = rest.indexOf(' ')
        if (secondSpace < 0) continue
        const path = rest.slice(0, secondSpace)
        const countsStr = rest.slice(secondSpace + 1).trim()

        const { unspecified, statusCodes, hitCount } = parseCounts(countsStr)

        routeList.push({
            method,
            path,
            tag: currentTag,
            hitCount,
            unspecified,
            statusCodes,
        })

        if (!serviceMap[currentTag]) {
            serviceMap[currentTag] = { name: currentTag, covered: 0, total: 0 }
        }
        serviceMap[currentTag].total++
        if (hitCount > 0) serviceMap[currentTag].covered++
    }

    const services = Object.values(serviceMap).map(s => ({
        ...s,
        coverage: s.total > 0 ? (s.covered / s.total) * 100 : 0,
    }))

    const totalRoutes = routeList.length
    const coveredRoutes = routeList.filter(r => r.hitCount > 0).length
    const coveragePercent = totalRoutes > 0 ? (coveredRoutes / totalRoutes) * 100 : 0

    return {
        totalRoutes,
        coveredRoutes,
        coveragePercent,
        routeList,
        services,
    }
}

function parseCounts(str) {
    const tokens = str.split(/\s+/)
    let unspecified = 0
    const statusCodes = {}

    for (const token of tokens) {
        if (token.includes(':')) {
            const [sc, cnt] = token.split(':')
            statusCodes[parseInt(sc, 10)] = parseInt(cnt, 10)
        } else {
            unspecified = parseInt(token, 10)
        }
    }

    let hitCount = unspecified
    for (const cnt of Object.values(statusCodes)) hitCount += cnt

    return { unspecified, statusCodes, hitCount }
}

export function routeCovClass(pct) {
    return pct >= config.thresholdGreen ? 'green' : pct >= config.thresholdYellow ? 'yellow' : 'red'
}

export function httpMethodClass(method) {
    switch (method) {
        case 'GET': return 'method-get'
        case 'POST': return 'method-post'
        case 'PUT': return 'method-put'
        case 'PATCH': return 'method-patch'
        case 'DELETE': return 'method-delete'
        default: return 'method-other'
    }
}

export function statusCodeClass(code) {
    if (code === -1) return 'status-fail'
    if (code >= 200 && code < 300) return 'status-2xx'
    if (code >= 300 && code < 400) return 'status-3xx'
    if (code >= 400 && code < 500) return 'status-4xx'
    if (code >= 500) return 'status-5xx'
    return 'status-other'
}
