import config from '@/config'

/**
 * Parse a Go coverprofile content string into structured data.
 */
export function parseCoverProfile(content) {
    const lines = content.trim().split('\n')
    if (!lines.length || !lines[0].startsWith('mode:')) {
        throw new Error('Invalid coverage profile: missing mode line')
    }
    const mode = lines[0].replace('mode:', '').trim()
    const blocks = []
    for (let i = 1; i < lines.length; i++) {
        const line = lines[i].trim()
        if (!line) continue
        const m = line.match(/^(.+):(\d+)\.(\d+),(\d+)\.(\d+)\s+(\d+)\s+(\d+)$/)
        if (m) {
            blocks.push({
                file: m[1],
                startLine: +m[2], startCol: +m[3],
                endLine: +m[4], endCol: +m[5],
                numStatements: +m[6], count: +m[7],
            })
        }
    }
    return { mode, blocks }
}

/**
 * Detect common module prefix from file paths.
 */
export function detectModulePrefix(filePaths) {
    if (!filePaths.length) return ''
    const parts = filePaths[0].split('/')
    let prefix = []
    for (let i = 0; i < parts.length - 1; i++) {
        const segment = parts[i]
        if (filePaths.every(p => p.split('/')[i] === segment)) {
            prefix.push(segment)
        } else break
    }
    return prefix.join('/')
}

/**
 * Build full coverage data model from parsed profile.
 */
export function buildCoverData(parsed) {
    const { mode, blocks } = parsed

    // Group blocks by file
    const fileMap = {}
    for (const b of blocks) {
        if (!fileMap[b.file]) fileMap[b.file] = { path: b.file, blocks: [], totalStatements: 0, coveredStatements: 0 }
        fileMap[b.file].blocks.push(b)
        fileMap[b.file].totalStatements += b.numStatements
        if (b.count > 0) fileMap[b.file].coveredStatements += b.numStatements
    }

    const filePaths = Object.keys(fileMap)
    const modulePrefix = detectModulePrefix(filePaths)

    // Compute per-file metrics
    const fileList = []
    for (const fp of filePaths) {
        const f = fileMap[fp]
        const rel = modulePrefix ? fp.slice(modulePrefix.length + 1) : fp
        const nameParts = rel.split('/')
        f.relativePath = rel
        f.name = nameParts[nameParts.length - 1]
        f.coverage = f.totalStatements > 0 ? (f.coveredStatements / f.totalStatements) * 100 : 0
        f.blocks.sort((a, b) => a.startLine - b.startLine || a.startCol - b.startCol)
        fileList.push(f)
    }

    // Build packages (group by directory)
    const pkgMap = {}
    for (const f of fileList) {
        const dir = f.relativePath.split('/').slice(0, -1).join('/') || '.'
        if (!pkgMap[dir]) pkgMap[dir] = { name: dir, shortName: dir.split('/').pop() || dir, files: [], totalStatements: 0, coveredStatements: 0 }
        pkgMap[dir].files.push(f)
        pkgMap[dir].totalStatements += f.totalStatements
        pkgMap[dir].coveredStatements += f.coveredStatements
    }
    const packageList = Object.values(pkgMap).map(p => {
        p.coverage = p.totalStatements > 0 ? (p.coveredStatements / p.totalStatements) * 100 : 0
        return p
    })

    // Totals
    let totalStatements = 0, totalCoveredStatements = 0
    for (const f of fileList) { totalStatements += f.totalStatements; totalCoveredStatements += f.coveredStatements }
    const totalCoverage = totalStatements > 0 ? (totalCoveredStatements / totalStatements) * 100 : 0

    // Build tree
    const tree = buildTree(fileList)

    return { mode, fileMap, fileList, packageList, totalStatements, totalCoveredStatements, totalCoverage, tree, modulePrefix }
}

/**
 * Build a folder/file tree from the flat file list.
 */
export function buildTree(fileList) {
    const root = { name: 'root', type: 'folder', children: [], allFiles: [] }
    for (const f of fileList) {
        const parts = f.relativePath.split('/')
        let node = root
        for (let i = 0; i < parts.length - 1; i++) {
            let child = node.children.find(c => c.type === 'folder' && c.name === parts[i])
            if (!child) {
                child = { name: parts[i], type: 'folder', children: [], allFiles: [], path: parts.slice(0, i + 1).join('/') }
                node.children.push(child)
            }
            node = child
        }
        const fileNode = { name: f.name, type: 'file', file: f, path: f.relativePath }
        node.children.push(fileNode)
    }
    computeFolderMetrics(root)
    sortTree(root)
    return root.children
}

function computeFolderMetrics(node) {
    if (node.type === 'file') return
    let totalStatements = 0, coveredStatements = 0
    const allFiles = []
    for (const child of node.children) {
        if (child.type === 'file') {
            totalStatements += child.file.totalStatements
            coveredStatements += child.file.coveredStatements
            allFiles.push(child.file)
        } else {
            computeFolderMetrics(child)
            totalStatements += child.totalStatements
            coveredStatements += child.coveredStatements
            allFiles.push(...child.allFiles)
        }
    }
    node.totalStatements = totalStatements
    node.coveredStatements = coveredStatements
    node.coverage = totalStatements > 0 ? (coveredStatements / totalStatements) * 100 : 0
    node.allFiles = allFiles
}

function sortTree(node) {
    if (node.type === 'file') return
    node.children.sort((a, b) => {
        if (a.type !== b.type) return a.type === 'folder' ? -1 : 1
        return a.name.localeCompare(b.name)
    })
    for (const child of node.children) sortTree(child)
}

/**
 * Filter tree nodes by search query.
 */
export function filterTreeNodes(nodes, query) {
    const result = []
    for (const node of nodes) {
        if (node.type === 'file') {
            if (node.name.toLowerCase().includes(query)) result.push(node)
        } else {
            const filteredChildren = filterTreeNodes(node.children || [], query)
            if (filteredChildren.length > 0 || node.name.toLowerCase().includes(query)) {
                result.push({ ...node, children: filteredChildren.length > 0 ? filteredChildren : node.children })
            }
        }
    }
    return result
}

/**
 * Find a tree node by path.
 */
export function findTreeNode(nodes, path) {
    for (const node of nodes) {
        if (node.path === path) return node
        if (node.type === 'folder' && node.children) {
            const found = findTreeNode(node.children, path)
            if (found) return found
        }
    }
    return null
}

/**
 * CSS class for coverage percentage.
 */
export function covClass(pct) {
    return pct >= config.thresholdGreen ? 'green' : pct >= config.thresholdYellow ? 'yellow' : 'red'
}

/**
 * Color string for coverage percentage.
 */
export function covColor(pct) {
    return pct >= config.thresholdGreen ? '#10b981' : pct >= config.thresholdYellow ? '#f59e0b' : '#ef4444'
}
