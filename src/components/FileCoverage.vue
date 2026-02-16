<template>
    <div class="file-cov-header">
        <h3><span class="material-icons-round">code</span> {{ file.name }}</h3>
        <div class="file-cov-meta">
            <span class="cov-badge cov-badge-lg" :class="covClass(file.coverage)">{{
                file.coverage.toFixed(1) }}%</span>
            <span class="text-sm text-muted">{{ file.coveredStatements }}/{{ file.totalStatements }} statements
                covered</span>
        </div>
    </div>

    <!-- Source code view if available -->
    <div v-if="sourceContent !== null" class="card file-cov-blocks source-view">
        <div v-for="(line, idx) in sourceLines" :key="idx" class="file-cov-line" :class="lineCovClass(idx + 1)">
            <div class="line-gutter">{{ idx + 1 }}</div>
            <div class="line-status"></div>
            <div class="line-code">{{ line }}</div>
        </div>
    </div>

    <!-- Block-based view if no source -->
    <div v-else>
        <div class="block-mode-toggle">
            <button class="btn btn-sm" :class="blockViewMode === 'visual' ? 'btn-primary' : 'btn-outline'"
                @click="blockViewMode = 'visual'">
                <span class="material-icons-round btn-icon">view_week</span> Visual Map
            </button>
            <button class="btn btn-sm" :class="blockViewMode === 'list' ? 'btn-primary' : 'btn-outline'"
                @click="blockViewMode = 'list'">
                <span class="material-icons-round btn-icon">list</span> Block List
            </button>
        </div>

        <!-- Visual line map -->
        <div v-if="blockViewMode === 'visual'" class="legend">
            <span class="legend-item"><span class="legend-swatch covered"></span>
                Covered</span>
            <span class="legend-item"><span class="legend-swatch uncovered"></span>
                Not covered</span>
            <span class="legend-item"><span class="legend-swatch neutral"></span>
                Not instrumented</span>
        </div>
        <div v-if="blockViewMode === 'visual'" class="card file-cov-blocks visual-map">
            <div v-for="line in lineMap" :key="line.num" class="file-cov-line" :class="line.cls">
                <div class="line-gutter">{{ line.num }}</div>
                <div class="line-status"></div>
                <div class="line-code visual-line-code">{{ line.label }}</div>
            </div>
        </div>

        <!-- Block list -->
        <div v-if="blockViewMode === 'list'" class="card block-list">
            <div v-for="(block, i) in file.blocks" :key="i" class="block-item"
                :class="block.count > 0 ? 'covered' : 'uncovered'">
                <div class="block-lines font-mono">L{{ block.startLine }}:{{ block.startCol }}–L{{ block.endLine }}:{{
                    block.endCol }}</div>
                <div class="block-stmts">{{ block.numStatements }} stmt{{ block.numStatements > 1 ? 's' : '' }}</div>
                <div class="block-hits" :style="{ color: block.count > 0 ? 'var(--green)' : 'var(--red)' }">
                    {{ block.count > 0 ? '✓ hit ' + block.count + 'x' : '✗ not hit' }}
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { covClass } from '@/utils/coverage'
import { getCodeFilesMap, hasInlineCodeFilesMap } from '@/utils/data'

export default {
    name: 'FileCoverage',
    props: {
        file: { type: Object, required: true },
        sourceRoot: { type: String, default: '' },
        modulePrefix: { type: String, default: '' },
    },
    data() {
        return {
            sourceContent: null,
            blockViewMode: 'visual',
        }
    },
    computed: {
        sourceLines() {
            if (!this.sourceContent) return []
            return this.sourceContent.split('\n')
        },
        lineMap() {
            const blocks = this.file.blocks
            if (!blocks.length) return []
            const maxLine = Math.max(...blocks.map(b => b.endLine))
            const lineStatus = new Array(maxLine + 1).fill(null)
            for (const b of blocks) {
                for (let l = b.startLine; l <= b.endLine; l++) {
                    if (b.count === 0) {
                        lineStatus[l] = { covered: false, count: 0, stmts: b.numStatements }
                    } else if (!lineStatus[l] || lineStatus[l].covered !== false) {
                        lineStatus[l] = { covered: true, count: b.count, stmts: b.numStatements }
                    }
                }
            }
            const result = []
            for (let i = 1; i <= maxLine; i++) {
                const s = lineStatus[i]
                if (s) {
                    result.push({
                        num: i,
                        cls: s.covered ? 'covered' : 'uncovered',
                        label: s.covered ? '✓ covered ×' + s.count : '✗ not covered',
                    })
                } else {
                    result.push({ num: i, cls: 'neutral', label: '' })
                }
            }
            return result
        },
    },
    watch: {
        file: {
            handler: 'tryFetchSource',
            immediate: true,
        },
    },
    methods: {
        covClass,
        lineCovClass(lineNum) {
            let result = null
            for (const b of this.file.blocks) {
                if (lineNum >= b.startLine && lineNum <= b.endLine) {
                    if (b.count === 0) return 'uncovered'
                    result = 'covered'
                }
            }
            return result || 'neutral'
        },
        async tryFetchSource() {
            this.sourceContent = null

            if (hasInlineCodeFilesMap()) {
                const map = await getCodeFilesMap()
                if (map) {
                    const content = map[this.file.path] || map[this.file.relativePath]
                    if (content) {
                        this.sourceContent = content
                        return
                    }
                }
            }

            if (!this.sourceRoot) return
            try {
                const relPath = this.modulePrefix
                    ? this.file.path.replace(this.modulePrefix + '/', '')
                    : this.file.relativePath
                const url = this.sourceRoot.replace(/\/$/, '') + '/' + relPath
                const resp = await fetch(url)
                if (resp.ok) {
                    this.sourceContent = await resp.text()
                }
            } catch (e) { /* source not available */ }
        },
    },
}
</script>

<style scoped>
.file-cov-meta {
    display: flex;
    gap: 16px;
    align-items: center;
    margin-top: 8px;
}

.cov-badge-lg {
    font-size: 14px;
    padding: 5px 14px;
}

.source-view {
    overflow: auto;
    max-height: calc(100vh - 240px);
}

.block-mode-toggle {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
}

.btn-icon {
    font-size: 15px;
}

.legend {
    display: flex;
    gap: 16px;
    margin-bottom: 12px;
    font-size: 12px;
    color: var(--text-secondary);
    align-items: center;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 4px;
}

.legend-swatch {
    width: 12px;
    height: 12px;
    border-radius: 2px;
    display: inline-block;
}

.legend-swatch.covered {
    background: var(--green);
}

.legend-swatch.uncovered {
    background: var(--red);
}

.legend-swatch.neutral {
    background: #e2e8f0;
}

.visual-map {
    overflow: auto;
    max-height: calc(100vh - 310px);
}

.visual-line-code {
    font-size: 12px;
}

.file-cov-blocks {
    font-family: 'JetBrains Mono', monospace;
    font-size: 13px;
    line-height: 1.7;
    border-radius: var(--radius);
    overflow: hidden;
    border: 1px solid var(--border);
}

.file-cov-line {
    display: flex;
}

.file-cov-line .line-gutter {
    width: 60px;
    text-align: right;
    padding: 0 10px;
    color: var(--text-muted);
    background: #f8fafc;
    border-right: 1px solid var(--border);
    user-select: none;
    flex-shrink: 0;
    font-size: 12px;
}

.file-cov-line .line-status {
    width: 4px;
    flex-shrink: 0;
}

.file-cov-line .line-code {
    flex: 1;
    padding: 0 14px;
    white-space: pre;
    overflow-x: auto;
}

.file-cov-line.covered {
    background: rgba(16, 185, 129, 0.08);
}

.file-cov-line.covered .line-status {
    background: var(--green);
}

.file-cov-line.uncovered {
    background: rgba(239, 68, 68, 0.08);
}

.file-cov-line.uncovered .line-status {
    background: var(--red);
}

.file-cov-line.neutral .line-status {
    background: transparent;
}

.block-list {
    border-radius: var(--radius);
    overflow: hidden;
    border: 1px solid var(--border);
}

.block-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 16px;
    border-bottom: 1px solid var(--border);
    font-size: 13px;
}

.block-item:last-child {
    border-bottom: none;
}

.block-item.covered {
    background: rgba(16, 185, 129, 0.04);
    border-left: 3px solid var(--green);
}

.block-item.uncovered {
    background: rgba(239, 68, 68, 0.04);
    border-left: 3px solid var(--red);
}

.block-item .block-lines {
    font-family: 'JetBrains Mono', monospace;
    font-weight: 500;
    min-width: 100px;
}

.block-item .block-stmts {
    color: var(--text-secondary);
    font-size: 12px;
}

.block-item .block-hits {
    margin-left: auto;
    font-weight: 600;
    font-size: 12px;
}
</style>
