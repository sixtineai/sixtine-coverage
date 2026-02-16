<template>
    <div class="content">
        <div class="card-grid">
            <div class="card stat-card" :class="covClass(coverData.totalCoverage)">
                <div class="stat-icon"><span class="material-icons-round">verified</span></div>
                <div class="stat-label">Total Coverage</div>
                <div class="stat-value">{{ coverData.totalCoverage.toFixed(1) }}%</div>
                <div class="stat-detail">{{ coverData.totalCoveredStatements }}/{{ coverData.totalStatements }}
                    statements</div>
            </div>
            <div class="card stat-card blue">
                <div class="stat-icon"><span class="material-icons-round">inventory_2</span></div>
                <div class="stat-label">Packages</div>
                <div class="stat-value">{{ coverData.packageList.length }}</div>
                <div class="stat-detail">Go packages analyzed</div>
            </div>
            <div class="card stat-card purple">
                <div class="stat-icon"><span class="material-icons-round">description</span></div>
                <div class="stat-label">Files</div>
                <div class="stat-value">{{ coverData.fileList.length }}</div>
                <div class="stat-detail">source files with coverage data</div>
            </div>
        </div>

        <div class="section-title"><span class="material-icons-round">apps</span> Coverage by Package</div>
        <div class="card-grid packages-grid">
            <div v-for="pkg in sortedPackages" :key="pkg.name" class="card pkg-card"
                @click="$emit('open-package', pkg)">
                <div class="pkg-name" :title="pkg.name">{{ pkg.shortName }}</div>
                <div class="pkg-cov" :style="{ color: covColor(pkg.coverage) }">{{ pkg.coverage.toFixed(1) }}%</div>
                <div class="progress-bar">
                    <div class="fill" :class="covClass(pkg.coverage)" :style="{ width: pkg.coverage + '%' }"></div>
                </div>
                <div class="pkg-detail">{{ pkg.coveredStatements }}/{{ pkg.totalStatements }} stmts &middot; {{
                    pkg.files.length }} files</div>
            </div>
        </div>
    </div>
</template>

<script>
import { covClass, covColor } from '@/utils/coverage'

export default {
    name: 'DashboardView',
    props: {
        coverData: { type: Object, required: true },
    },
    emits: ['open-package'],
    computed: {
        sortedPackages() {
            return [...this.coverData.packageList].sort((a, b) => a.name.localeCompare(b.name))
        },
    },
    methods: {
        covClass,
        covColor,
    },
}
</script>

<style scoped>
.packages-grid {
    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
}
</style>
