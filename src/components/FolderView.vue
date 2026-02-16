<template>
    <div v-if="folder">
        <div class="file-cov-header">
            <h3><span class="material-icons-round">folder</span> {{ folder.name }}</h3>
            <div class="text-sm text-muted">{{ folder.allFiles.length }} files &middot; {{ folder.totalStatements }}
                statements</div>
        </div>
        <div class="card-grid folder-stats">
            <div class="card stat-card" :class="covClass(folder.coverage)">
                <div class="stat-icon"><span class="material-icons-round">verified</span></div>
                <div class="stat-label">Folder Coverage</div>
                <div class="stat-value">{{ folder.coverage.toFixed(1) }}%</div>
                <div class="stat-detail">{{ folder.coveredStatements }}/{{ folder.totalStatements }} statements</div>
            </div>
            <div class="card stat-card blue">
                <div class="stat-icon"><span class="material-icons-round">description</span></div>
                <div class="stat-label">Files</div>
                <div class="stat-value">{{ folder.allFiles.length }}</div>
                <div class="stat-detail">in this directory</div>
            </div>
        </div>
        <div class="section-title"><span class="material-icons-round">list</span> Files in folder</div>
        <div class="card table-card">
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th class="col-coverage">Coverage</th>
                        <th class="col-progress">Progress</th>
                        <th class="col-statements">Statements</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="f in folder.allFiles" :key="f.path" class="clickable-row"
                        @click="$emit('open-file', f)">
                        <td>
                            <div class="file-name"><span class="material-icons-round">description</span> {{ f.name }}
                            </div>
                        </td>
                        <td><span class="cov-badge" :class="covClass(f.coverage)">{{ f.coverage.toFixed(1) }}%</span>
                        </td>
                        <td class="progress-cell">
                            <div class="progress-bar">
                                <div class="fill" :class="covClass(f.coverage)" :style="{ width: f.coverage + '%' }">
                                </div>
                            </div>
                        </td>
                        <td class="stmt-cell">{{ f.coveredStatements }}/{{ f.totalStatements }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    <div v-else class="empty-state">
        <span class="material-icons-round">account_tree</span>
        <h3>Select a folder or file</h3>
        <p class="text-sm text-muted">Navigate the tree on the left to view coverage details</p>
    </div>
</template>

<script>
import { covClass } from '@/utils/coverage'

export default {
    name: 'FolderView',
    props: {
        folder: { type: Object, default: null },
    },
    emits: ['open-file'],
    methods: {
        covClass,
    },
}
</script>

<style scoped>
.table-card {
    overflow: hidden;
}

.folder-stats {
    margin-bottom: 20px;
}

.clickable-row {
    cursor: pointer;
}

.col-coverage {
    width: 100px;
}

.col-progress {
    width: 180px;
}

.col-statements {
    width: 130px;
}
</style>
