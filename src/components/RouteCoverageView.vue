<template>
    <div class="content">
        <!-- Dashboard cards -->
        <div class="card-grid">
            <div class="card stat-card" :class="routeCovClass(routeData.coveragePercent)">
                <div class="stat-icon"><span class="material-icons-round">alt_route</span></div>
                <div class="stat-label">Route Coverage</div>
                <div class="stat-value">{{ routeData.coveragePercent.toFixed(1) }}%</div>
                <div class="stat-detail">{{ routeData.coveredRoutes }}/{{ routeData.totalRoutes }} routes covered</div>
            </div>
            <div class="card stat-card blue">
                <div class="stat-icon"><span class="material-icons-round">send</span></div>
                <div class="stat-label">Total Calls</div>
                <div class="stat-value">{{ totalCalls }}</div>
                <div class="stat-detail">HTTP requests recorded</div>
            </div>
            <div class="card stat-card purple">
                <div class="stat-icon"><span class="material-icons-round">dns</span></div>
                <div class="stat-label">Services</div>
                <div class="stat-value">{{ routeData.services.length }}</div>
                <div class="stat-detail">API service tags</div>
            </div>
            <div class="card stat-card" :class="uncoveredRoutes.length === 0 ? 'green' : 'red'">
                <div class="stat-icon"><span class="material-icons-round">report_off</span></div>
                <div class="stat-label">Uncovered</div>
                <div class="stat-value">{{ uncoveredRoutes.length }}</div>
                <div class="stat-detail">routes with no test calls</div>
            </div>
        </div>

        <!-- Service breakdown -->
        <div class="section-title"><span class="material-icons-round">apps</span> Coverage by Service</div>
        <div class="card-grid service-grid">
            <div v-for="svc in sortedServices" :key="svc.name" class="card pkg-card">
                <div class="pkg-name" :title="svc.name">{{ svc.name }}</div>
                <div class="pkg-cov" :style="{ color: covColor(svc.coverage) }">{{ svc.coverage.toFixed(1) }}%</div>
                <div class="progress-bar">
                    <div class="fill" :class="routeCovClass(svc.coverage)" :style="{ width: svc.coverage + '%' }">
                    </div>
                </div>
                <div class="pkg-detail">{{ svc.covered }}/{{ svc.total }} routes</div>
            </div>
        </div>

        <!-- Routes table -->
        <div class="section-title"><span class="material-icons-round">list</span> Covered Routes</div>
        <div class="search-bar">
            <div class="search-wrapper">
                <span class="material-icons-round search-icon">search</span>
                <input v-model="search" placeholder="Search routes..." class="search-input">
            </div>
            <select v-model="methodFilter" class="route-select">
                <option value="">All Methods</option>
                <option v-for="m in availableMethods" :key="m" :value="m">{{ m }}</option>
            </select>
        </div>
        <div class="card table-card">
            <table class="data-table">
                <thead>
                    <tr>
                        <th @click="setSort('method')" class="col-method">Method <span v-if="sortField === 'method'"
                                class="sort-icon">{{ sortDir === 1 ? '↑' : '↓' }}</span></th>
                        <th @click="setSort('path')">Path <span v-if="sortField === 'path'" class="sort-icon">{{
                            sortDir === 1 ? '↑' : '↓' }}</span></th>
                        <th @click="setSort('tag')" class="col-service">Service <span
                                v-if="sortField === 'tag'" class="sort-icon">{{ sortDir === 1 ? '↑' : '↓'
                                }}</span></th>
                        <th @click="setSort('hitCount')" class="col-hits">Hits <span
                                v-if="sortField === 'hitCount'" class="sort-icon">{{ sortDir === 1 ? '↑' : '↓'
                                }}</span></th>
                        <th class="col-status">Status Codes</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="route in filteredRoutes" :key="route.method + route.path"
                        :class="{ 'uncovered-row': route.hitCount === 0 }">
                        <td><span class="http-method-badge" :class="httpMethodClass(route.method)">{{ route.method
                                }}</span></td>
                        <td class="font-mono text-sm">{{ route.path }}</td>
                        <td><span class="service-tag">{{ route.tag }}</span></td>
                        <td class="hits-cell">
                            <span class="hit-count" :class="{ 'hit-zero': route.hitCount === 0 }">{{ route.hitCount }}</span>
                        </td>
                        <td>
                            <div class="status-badges" v-if="route.hitCount > 0">
                                <span v-if="route.unspecified > 0"
                                    class="status-badge status-other">
                                    ??? <span class="status-count">×{{ route.unspecified }}</span>
                                </span>
                                <span v-for="(count, code) in route.statusCodes" :key="code"
                                    class="status-badge" :class="statusCodeClass(Number(code))">
                                    {{ Number(code) === -1 ? 'FAIL' : code }} <span class="status-count">×{{ count }}</span>
                                </span>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div v-if="filteredRoutes.length === 0" class="empty-state">
                <span class="material-icons-round">search_off</span>
                <h3>No routes match your filter</h3>
            </div>
        </div>

        <!-- Uncovered routes -->
        <div v-if="uncoveredRoutes.length > 0" class="uncovered-section">
            <div class="section-title"><span class="material-icons-round">report_off</span> Uncovered Routes</div>
            <div class="card table-card">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th class="col-method">Method</th>
                            <th>Path</th>
                            <th class="col-service-sm">Service</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="route in uncoveredRoutes"
                            :key="route.method + route.path" class="uncovered-row">
                            <td><span class="http-method-badge" :class="httpMethodClass(route.method)">{{ route.method
                                    }}</span></td>
                            <td class="font-mono text-sm">{{ route.path }}</td>
                            <td><span class="service-tag">{{ route.tag }}</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script>
import { covColor } from '@/utils/coverage'
import { routeCovClass, httpMethodClass, statusCodeClass } from '@/utils/routeCoverage'

export default {
    name: 'RouteCoverageView',
    props: {
        routeData: { type: Object, required: true },
    },
    data() {
        return {
            search: '',
            methodFilter: '',
            sortField: 'hitCount',
            sortDir: -1,
        }
    },
    computed: {
        totalCalls() {
            return this.routeData.routeList.reduce((sum, r) => sum + r.hitCount, 0)
        },
        uncoveredRoutes() {
            return this.routeData.routeList.filter(r => r.hitCount === 0)
        },
        sortedServices() {
            return [...this.routeData.services].sort((a, b) => a.name.localeCompare(b.name))
        },
        availableMethods() {
            const methods = new Set(this.routeData.routeList.map(r => r.method))
            return [...methods].sort()
        },
        filteredRoutes() {
            let list = [...this.routeData.routeList]
            if (this.search) {
                const q = this.search.toLowerCase()
                list = list.filter(r =>
                    r.path.toLowerCase().includes(q) ||
                    r.tag.toLowerCase().includes(q) ||
                    r.method.toLowerCase().includes(q)
                )
            }
            if (this.methodFilter) {
                list = list.filter(r => r.method === this.methodFilter)
            }
            const field = this.sortField
            const dir = this.sortDir
            list.sort((a, b) => {
                if (field === 'method') return a.method.localeCompare(b.method) * dir
                if (field === 'path') return a.path.localeCompare(b.path) * dir
                if (field === 'tag') return a.tag.localeCompare(b.tag) * dir
                if (field === 'hitCount') return (a.hitCount - b.hitCount) * dir
                return 0
            })
            return list
        },
    },
    methods: {
        covColor,
        routeCovClass,
        httpMethodClass,
        statusCodeClass,
        setSort(field) {
            if (this.sortField === field) {
                this.sortDir *= -1
            } else {
                this.sortField = field
                this.sortDir = field === 'hitCount' ? -1 : 1
            }
        },
    },
}
</script>

<style scoped>
.service-grid {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
    margin-bottom: 28px;
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

.route-select {
    padding: 8px 12px;
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

.hits-cell {
    text-align: center;
}

.uncovered-section {
    margin-top: 28px;
}

.col-method {
    width: 80px;
}

.col-service {
    width: 180px;
}

.col-hits {
    width: 80px;
}

.col-status {
    width: 220px;
}

.col-service-sm {
    width: 140px;
}

.http-method-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 11px;
    font-weight: 700;
    font-family: 'JetBrains Mono', monospace;
    min-width: 52px;
    text-align: center;
}

.method-get {
    background: #dbeafe;
    color: #1d4ed8;
}

.method-post {
    background: #dcfce7;
    color: #15803d;
}

.method-put {
    background: #fef3c7;
    color: #b45309;
}

.method-patch {
    background: #fae8ff;
    color: #a21caf;
}

.method-delete {
    background: #fef2f2;
    color: #dc2626;
}

.method-other {
    background: #f1f5f9;
    color: #475569;
}

.status-badges {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.status-badge {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 2px 8px;
    border-radius: 99px;
    font-size: 11px;
    font-weight: 600;
    font-family: 'JetBrains Mono', monospace;
}

.status-count {
    font-weight: 400;
    opacity: 0.7;
    font-size: 10px;
}

.status-2xx {
    background: var(--green-bg);
    color: #065f46;
    border: 1px solid var(--green-border);
}

.status-3xx {
    background: #dbeafe;
    color: #1e40af;
    border: 1px solid #bfdbfe;
}

.status-4xx {
    background: var(--yellow-bg);
    color: #92400e;
    border: 1px solid var(--yellow-border);
}

.status-5xx {
    background: var(--red-bg);
    color: #991b1b;
    border: 1px solid var(--red-border);
}

.status-fail {
    background: #fef2f2;
    color: #991b1b;
    border: 1px solid #fecaca;
}

.status-other {
    background: #f1f5f9;
    color: #64748b;
    border: 1px solid #e2e8f0;
}

.hit-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: #f1f5f9;
    color: var(--text-primary);
    padding: 3px 10px;
    border-radius: 99px;
    font-size: 12px;
    font-weight: 700;
    font-family: 'JetBrains Mono', monospace;
    min-width: 32px;
}

.hit-zero {
    opacity: 0.4;
}

.service-tag {
    font-size: 12px;
    color: var(--text-secondary);
    background: #f1f5f9;
    padding: 3px 8px;
    border-radius: 4px;
    font-weight: 500;
}

.uncovered-row td {
    opacity: 0.75;
}

@media (max-width: 768px) {
    .table-card {
        overflow-x: auto;
        -webkit-overflow-scrolling: touch;
    }

    .col-service,
    .col-status {
        display: none;
    }

    .data-table td:nth-child(3),
    .data-table th:nth-child(3),
    .data-table td:nth-child(5),
    .data-table th:nth-child(5) {
        display: none;
    }

    .service-grid {
        grid-template-columns: 1fr;
    }

    .search-bar {
        flex-wrap: wrap;
    }

    .route-select {
        width: 100%;
    }
}
</style>
