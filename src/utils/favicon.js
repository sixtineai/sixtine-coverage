import config from '@/config'

function faviconColor(pct) {
    if (pct >= config.thresholdGreen) return '#059669'
    if (pct >= config.thresholdYellow) return '#d97706'
    return '#dc2626'
}

function generateFaviconSVG(percentage) {
    const color = percentage != null ? faviconColor(percentage) : '#0d9488'
    if (percentage == null) {
        return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="${color}"/>
  <text x="16" y="17" text-anchor="middle" dominant-baseline="central"
    font-family="Inter,system-ui,sans-serif" font-weight="700"
    font-size="14" fill="#fff">Go</text>
</svg>`
    }
    const num = Math.round(percentage)
    return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect width="32" height="32" rx="6" fill="${color}"/>
  <text x="16" y="16" text-anchor="middle" dominant-baseline="central"
    font-family="Inter,system-ui,sans-serif" font-weight="800"
    font-size="24" fill="#fff">${num}</text>
</svg>`
}

export function setFavicon(percentage) {
    const svg = generateFaviconSVG(percentage)
    const encoded = 'data:image/svg+xml,' + encodeURIComponent(svg)
    let link = document.querySelector('link[rel="icon"]')
    if (!link) {
        link = document.createElement('link')
        link.rel = 'icon'
        document.head.appendChild(link)
    }
    link.type = 'image/svg+xml'
    link.href = encoded
}
