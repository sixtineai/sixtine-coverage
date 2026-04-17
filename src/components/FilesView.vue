<template>
    <div class="content">
        <div class="search-bar">
            <div class="search-wrapper">
                <span class="material-icons-round search-icon">search</span>
                <input v-model="fileSearch" placeholder="Search files..." class="search-input">
            </div>
            <button class="btn btn-outline btn-sm" @click="toggleSort">
                <span class="material-icons-round sort-btn-icon">sort</span>
                Sort: {{ sortLabel }}
            </button>
        </div>
        <div class="card table-card">
            <table class="data-table">
                <thead>
                    <tr>
                        <th @click="setSort('name')">Name <span v-if="sortField === 'name'" class="sort-icon">{{ sortDir
                                === 1 ? '↑' : '↓' }}</span></th>
                        <th @click="setSort('coverage')" class="col-coverage">Coverage <span
                                v-if="sortField === 'coverage'" class="sort-icon">{{ sortDir === 1 ? '↑' : '↓' }}</span>
                        </th>
                        <th class="col-progress">Progress</th>
                        <th @click="setSort('statements')" class="col-statements">Statements <span
                                v-if="sortField === 'statements'" class="sort-icon">{{ sortDir === 1 ? '↑' : '↓' }}</span>
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="file in filteredFiles" :key="file.path" class="clickable-row"
                        @click="$emit('open-file', file)">
                        <td>
                            <div class="file-name">
                                <span class="material-icons-round">description</span>
                                <span :title="file.path">{{ file.relativePath }}</span>
                            </div>
                        </td>
                        <td><span class="cov-badge" :class="covClass(file.coverage)">{{ file.coverage.toFixed(1)
                                }}%</span></td>
                        <td class="progress-cell">
                            <div class="progress-bar">
                                <div class="fill" :class="covClass(file.coverage)"
                                    :style="{ width: file.coverage + '%' }"></div>
                            </div>
                        </td>
                        <td class="stmt-cell">{{ file.coveredStatements }}/{{ file.totalStatements }}</td>
                    </tr>
                </tbody>
            </table>
            <div v-if="filteredFiles.length === 0" class="empty-state">
                <span class="material-icons-round">search_off</span>
                <h3>No files match your search</h3>
            </div>
        </div>
    </div>
</template>

<script>
import { covClass } from '@/utils/coverage'

export default {
    name: 'FilesView',
    props: {
        files: { type: Array, required: true },
    },
    emits: ['open-file'],
    data() {
        return {
            fileSearch: '',
            sortField: 'name',
            sortDir: 1,
        }
    },
    computed: {
        sortLabel() {
            return { name: 'Name', coverage: 'Coverage', statements: 'Statements' }[this.sortField] || 'Name'
        },
        filteredFiles() {
            let list = [...this.files]
            if (this.fileSearch) {
                const q = this.fileSearch.toLowerCase()
                list = list.filter(f => f.relativePath.toLowerCase().includes(q))
            }
            const sort = this.sortField
            const dir = this.sortDir
            list.sort((a, b) => {
                if (sort === 'name') return a.relativePath.localeCompare(b.relativePath) * dir
                if (sort === 'coverage') return (a.coverage - b.coverage) * dir
                if (sort === 'statements') return (a.totalStatements - b.totalStatements) * dir
                return 0
            })
            return list
        },
    },
    methods: {
        covClass,
        setSort(field) {
            if (this.sortField === field) {
                this.sortDir *= -1
            } else {
                this.sortField = field
                this.sortDir = field === 'name' ? 1 : -1
            }
        },
        toggleSort() {
            const fields = ['name', 'coverage', 'statements']
            const idx = fields.indexOf(this.sortField)
            this.sortField = fields[(idx + 1) % fields.length]
            this.sortDir = this.sortField === 'name' ? 1 : -1
        },
    },
}
</script>

<style scoped>
.sort-btn-icon {
    font-size: 16px;
}

.search-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
    align-items: center;
}

.search-wrapper {
    flex: 1;
    position: relative;
}

.search-icon {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 18px;
    color: var(--text-muted);
}

.search-input {
    width: 100%;
    padding: 9px 12px 9px 36px;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 13px;
    font-family: inherit;
    outline: none;
    background: var(--card-bg);
}

.table-card {
    overflow: hidden;
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
    width: 140px;
}

@media (max-width: 768px) {
    .table-card {
        overflow-x: auto;
        -webkit-overflow-scrolling: touch;
    }

    .col-progress,
    .col-statements {
        display: none;
    }

    .data-table td:nth-child(3),
    .data-table th:nth-child(3),
    .data-table td:nth-child(4),
    .data-table th:nth-child(4) {
        display: none;
    }

    .search-bar {
        flex-wrap: wrap;
    }
}
</style>
