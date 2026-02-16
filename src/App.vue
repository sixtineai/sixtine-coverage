<template>
    <!-- Loading / Config Screen -->
    <ConfigScreen v-if="!loaded" :initial-profile-url="profileUrl" :initial-source-root="sourceRoot"
        :initial-module-prefix="modulePrefix" :initial-route-coverage-url="routeCoverageUrl" @loaded="onLoaded" />

    <!-- Main App -->
    <template v-else>
        <!-- Sidebar -->
        <aside class="sidebar" :class="{ open: sidebarOpen }">
            <div class="sidebar-header">
                <div class="logo">Go</div>
                <div>
                    <h1>Coverage</h1>
                    <span>{{ coverData.mode }} mode</span>
                </div>
            </div>
            <nav class="sidebar-nav">
                <div class="nav-section">
                    <div class="nav-section-title">Navigation</div>
                    <div class="nav-item" :class="{ active: view === 'dashboard' }" @click="navigate('dashboard')">
                        <span class="material-icons-round">dashboard</span> Dashboard
                    </div>
                    <div class="nav-item" :class="{ active: view === 'files' }" @click="navigate('files')">
                        <span class="material-icons-round">description</span> Files
                        <span class="badge">{{ coverData.fileList.length }}</span>
                    </div>
                    <div class="nav-item" :class="{ active: isExplorerActive }" @click="navigate('tree')">
                        <span class="material-icons-round">account_tree</span> Explorer
                    </div>
                    <div v-if="hasRouteCoverage" class="nav-item" :class="{ active: view === 'routes' }" @click="navigate('routes')">
                        <span class="material-icons-round">alt_route</span> Routes
                        <span class="badge">{{ routeData.coveredRoutes }}/{{ routeData.totalRoutes }}</span>
                    </div>
                </div>
                <div class="nav-section">
                    <div class="nav-section-title">Summary</div>
                    <div class="summary-wrapper">
                        <div class="summary-row">
                            <span class="summary-label">Total Coverage</span>
                            <span class="summary-value"
                                :style="{ color: covColor(coverData.totalCoverage) }">{{
                                    coverData.totalCoverage.toFixed(1) }}%</span>
                        </div>
                        <div class="progress-bar summary-progress">
                            <div class="fill" :class="covClass(coverData.totalCoverage)"
                                :style="{ width: coverData.totalCoverage + '%' }"></div>
                        </div>
                        <div class="summary-detail">
                            {{ coverData.totalCoveredStatements }}/{{ coverData.totalStatements }} statements
                        </div>
                        <template v-if="routeData">
                            <div class="summary-row">
                                <span class="summary-label">Route Coverage</span>
                                <span class="summary-value"
                                    :style="{ color: covColor(routeData.coveragePercent) }">{{
                                        routeData.coveragePercent.toFixed(1) }}%</span>
                            </div>
                            <div class="progress-bar summary-progress">
                                <div class="fill" :class="covClass(routeData.coveragePercent)"
                                    :style="{ width: routeData.coveragePercent + '%' }"></div>
                            </div>
                            <div class="summary-detail">
                                {{ routeData.coveredRoutes }}/{{ routeData.totalRoutes }} routes
                            </div>
                        </template>
                    </div>
                </div>
                <div class="nav-section">
                    <div class="nav-section-title">Config</div>
                    <div class="nav-item" @click="loaded = false">
                        <span class="material-icons-round">settings</span> Change Profile
                    </div>
                </div>
            </nav>
            <div class="sidebar-footer">
                Go Coverage Viewer &middot; Made by SixtineAI
            </div>
        </aside>

        <!-- Mobile overlay -->
        <div class="sidebar-overlay" :class="{ visible: sidebarOpen }" @click="sidebarOpen = false"></div>

        <!-- Main Content -->
        <div class="main">
            <div class="topbar">
                <div class="topbar-left">
                    <button class="burger-btn" @click="sidebarOpen = !sidebarOpen">
                        <span class="material-icons-round">menu</span>
                    </button>
                    <h2>
                        <span class="material-icons-round">{{ viewIcon }}</span>
                        {{ viewTitle }}
                    </h2>
                </div>
                <div class="breadcrumb" v-if="breadcrumb.length > 0">
                    <a @click="navigate('dashboard')">Home</a>
                    <template v-for="(b, i) in breadcrumb" :key="i">
                        <span>/</span>
                        <a v-if="i < breadcrumb.length - 1" @click="b.action">{{ b.label }}</a>
                        <span v-else class="breadcrumb-current">{{ b.label }}</span>
                    </template>
                </div>
            </div>

            <!-- Dashboard View -->
            <DashboardView v-if="view === 'dashboard'" :cover-data="coverData" @open-package="openPackage" />

            <!-- Files View -->
            <FilesView v-if="view === 'files'" :files="coverData.fileList" @open-file="openFile" />

            <!-- Tree / Explorer View -->
            <TreeExplorer v-if="isExplorerActive" ref="treeExplorer" :tree="coverData.tree" :source-root="sourceRoot"
                :module-prefix="modulePrefix" :initial-folder-path="explorerInitPath"
                @view-change="onExplorerViewChange" />

            <!-- Route Coverage View -->
            <RouteCoverageView v-if="view === 'routes' && routeData" :route-data="routeData" />
        </div>
    </template>
</template>

<script>
import { covClass, covColor, findTreeNode, parseCoverProfile, buildCoverData } from '@/utils/coverage'
import { parseRouteCoverage } from '@/utils/routeCoverage'
import { resolveInlineCoverage, resolveInlineRouteCoverage } from '@/utils/data'
import { setFavicon } from '@/utils/favicon'
import config from '@/config'
import ConfigScreen from '@/components/ConfigScreen.vue'
import DashboardView from '@/components/DashboardView.vue'
import FilesView from '@/components/FilesView.vue'
import TreeExplorer from '@/components/TreeExplorer.vue'
import RouteCoverageView from '@/components/RouteCoverageView.vue'

export default {
    name: 'App',
    components: { ConfigScreen, DashboardView, FilesView, TreeExplorer, RouteCoverageView },
    data() {
        return {
            loaded: false,
            sidebarOpen: false,
            profileUrl: config.coverageUrl,
            sourceRoot: config.sourceRoot,
            modulePrefix: config.modulePrefix,
            routeCoverageUrl: config.routeCoverageUrl,
            coverData: null,
            routeData: null,
            view: 'dashboard',
            explorerInitPath: '',
            explorerSubView: 'tree',
            explorerSelectedFile: null,
            explorerSelectedFolder: null,
        }
    },
    computed: {
        autoLoad() {
            return !!config.coverageData || (!!config.coverageUrl && config.coverageUrl !== './cover.out')
        },
        hasRouteCoverage() {
            return !!this.routeData
        },
        isExplorerActive() {
            return this.view === 'tree' || this.view === 'folder' || this.view === 'filecov'
        },
        viewTitle() {
            switch (this.view) {
                case 'dashboard': return 'Dashboard'
                case 'files': return 'All Files'
                case 'tree': return 'Explorer'
                case 'folder': return this.explorerSelectedFolder ? this.explorerSelectedFolder.name : 'Explorer'
                case 'filecov': return this.explorerSelectedFile ? this.explorerSelectedFile.relativePath : 'File Coverage'
                case 'routes': return 'Route Coverage'
                default: return 'Dashboard'
            }
        },
        viewIcon() {
            switch (this.view) {
                case 'dashboard': return 'dashboard'
                case 'files': return 'description'
                case 'tree': case 'folder': return 'account_tree'
                case 'filecov': return 'code'
                case 'routes': return 'alt_route'
                default: return 'dashboard'
            }
        },
        breadcrumb() {
            if (this.view === 'dashboard') return []
            const crumbs = []
            if (this.view === 'files') {
                crumbs.push({ label: 'Files' })
            } else if (this.view === 'routes') {
                crumbs.push({ label: 'Route Coverage' })
            } else if (this.view === 'tree') {
                crumbs.push({ label: 'Explorer' })
            } else if (this.view === 'folder' && this.explorerSelectedFolder) {
                crumbs.push({ label: 'Explorer', action: () => this.navigate('tree') })
                const parts = this.explorerSelectedFolder.path
                    ? this.explorerSelectedFolder.path.split('/')
                    : [this.explorerSelectedFolder.name]
                parts.forEach(p => {
                    crumbs.push({ label: p })
                })
            } else if (this.view === 'filecov' && this.explorerSelectedFile) {
                crumbs.push({ label: 'Explorer', action: () => this.navigate('tree') })
                crumbs.push({ label: this.explorerSelectedFile.name })
            }
            return crumbs
        },
    },
    mounted() {
        setFavicon(null)
        window.addEventListener('popstate', this.onPopState)
        if (this.autoLoad) {
            this.$nextTick(() => {
                this.autoLoadProfile()
            })
        }
        if (!config.coverageData) {
            if (this.routeCoverageUrl) {
                this.loadRouteCoverage()
            }
        }
    },
    beforeUnmount() {
        window.removeEventListener('popstate', this.onPopState)
    },
    methods: {
        covClass,
        covColor,
        async autoLoadProfile() {
            try {
                let text = await resolveInlineCoverage()
                if (!text) {
                    const resp = await fetch(this.profileUrl)
                    if (!resp.ok) throw new Error('HTTP ' + resp.status)
                    text = await resp.text()
                }
                const parsed = parseCoverProfile(text)
                const coverData = buildCoverData(parsed)
                if (!this.modulePrefix) this.modulePrefix = coverData.modulePrefix
                this.coverData = coverData
                this.loaded = true
                this.view = 'dashboard'
                setFavicon(coverData.totalCoverage)
                history.replaceState(this.getViewState(), '', '')

                const inlineRouteData = await resolveInlineRouteCoverage()
                if (inlineRouteData) {
                    this.routeData = parseRouteCoverage(inlineRouteData)
                } else if (this.routeCoverageUrl) {
                    await this.loadRouteCoverage()
                }
            } catch (e) {
                console.error('Auto-load failed:', e)
            }
        },
        async loadRouteCoverage() {
            try {
                const resp = await fetch(this.routeCoverageUrl)
                if (!resp.ok) throw new Error('HTTP ' + resp.status)
                const text = await resp.text()
                this.routeData = parseRouteCoverage(text)
            } catch (e) {
                console.error('Route coverage load failed:', e)
            }
        },
        onLoaded({ coverData, profileUrl, sourceRoot, modulePrefix, routeCoverageUrl }) {
            this.coverData = coverData
            this.profileUrl = profileUrl
            this.sourceRoot = sourceRoot
            this.modulePrefix = modulePrefix
            setFavicon(coverData.totalCoverage)
            if (routeCoverageUrl) {
                this.routeCoverageUrl = routeCoverageUrl
                this.loadRouteCoverage()
            }
            this.loaded = true
            this.view = 'dashboard'
            history.replaceState(this.getViewState(), '', '')
        },
        navigate(view) {
            this.sidebarOpen = false
            this.view = view
            this.explorerInitPath = ''
            this.explorerSelectedFile = null
            this.explorerSelectedFolder = null
            this.pushHistory()
        },
        openPackage(pkg) {
            this.view = 'tree'
            this.explorerInitPath = pkg.name
            this.pushHistory()
        },
        openFile(file) {
            this.explorerSelectedFile = file
            this.view = 'filecov'
            this.explorerInitPath = file.relativePath
            this.pushHistory()
        },
        onExplorerViewChange({ view, node, file }) {
            this.view = view
            if (view === 'folder') {
                this.explorerSelectedFolder = node
                this.explorerSelectedFile = null
            } else if (view === 'filecov') {
                this.explorerSelectedFile = file
            }
            this.pushHistory()
        },
        getViewState() {
            return {
                view: this.view,
                explorerInitPath: this.explorerInitPath,
                filePath: this.explorerSelectedFile?.relativePath || '',
                folderPath: this.explorerSelectedFolder?.path || '',
            }
        },
        pushHistory() {
            history.pushState(this.getViewState(), '', '')
        },
        onPopState(event) {
            if (!this.loaded) return
            const state = event.state
            if (!state) {
                this.view = 'dashboard'
                this.explorerInitPath = ''
                this.explorerSelectedFile = null
                this.explorerSelectedFolder = null
                return
            }
            this.view = state.view || 'dashboard'
            this.explorerInitPath = state.explorerInitPath || ''
            if (state.filePath && this.coverData) {
                this.explorerSelectedFile = this.coverData.fileList.find(f => f.relativePath === state.filePath) || null
            } else {
                this.explorerSelectedFile = null
            }
            if (state.folderPath && this.coverData) {
                this.explorerSelectedFolder = findTreeNode(this.coverData.tree, state.folderPath) || null
            } else {
                this.explorerSelectedFolder = null
            }
        },
    },
}
</script>

<style scoped>
.sidebar {
    width: var(--sidebar-width);
    background: var(--sidebar-bg);
    color: var(--sidebar-text);
    display: flex;
    flex-direction: column;
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 100;
    transition: transform 0.3s ease;
}

.sidebar-header {
    padding: 20px 18px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    display: flex;
    align-items: center;
    gap: 10px;
}

.sidebar-header .logo {
    width: 32px;
    height: 32px;
    background: var(--teal);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 700;
    font-size: 14px;
}

.sidebar-header h1 {
    font-size: 15px;
    font-weight: 600;
    color: var(--sidebar-text-active);
    letter-spacing: -0.3px;
}

.sidebar-header span {
    font-size: 11px;
    color: var(--sidebar-text);
    display: block;
    margin-top: 1px;
}

.sidebar-nav {
    flex: 1;
    padding: 12px 10px;
    overflow-y: auto;
}

.nav-section {
    margin-bottom: 20px;
}

.nav-section-title {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 1.2px;
    color: rgba(148, 163, 184, 0.5);
    padding: 0 10px;
    margin-bottom: 6px;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 9px 12px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 13.5px;
    font-weight: 450;
    color: var(--sidebar-text);
    transition: all 0.15s ease;
    text-decoration: none;
    user-select: none;
}

.nav-item:hover {
    background: var(--sidebar-hover);
    color: var(--sidebar-text-active);
}

.nav-item.active {
    background: var(--teal);
    color: #fff;
    font-weight: 500;
}

.nav-item .material-icons-round {
    font-size: 19px;
}

.nav-item .badge {
    margin-left: auto;
    background: rgba(255, 255, 255, 0.1);
    color: var(--sidebar-text);
    font-size: 11px;
    padding: 2px 7px;
    border-radius: 10px;
    font-weight: 500;
}

.nav-item.active .badge {
    background: rgba(255, 255, 255, 0.25);
    color: #fff;
}

.sidebar-footer {
    padding: 14px 18px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
    font-size: 11px;
    color: rgba(148, 163, 184, 0.5);
}

.main {
    margin-left: var(--sidebar-width);
    flex: 1;
    min-height: 100vh;
}

.topbar {
    background: var(--card-bg);
    border-bottom: 1px solid var(--border);
    padding: 14px 28px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: sticky;
    top: 0;
    z-index: 50;
}

.topbar h2 {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 8px;
}

.topbar h2 .material-icons-round {
    font-size: 20px;
    color: var(--teal);
}

.breadcrumb {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--text-muted);
}

.breadcrumb a {
    color: var(--teal);
    text-decoration: none;
    cursor: pointer;
}

.breadcrumb a:hover {
    text-decoration: underline;
}

.summary-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
}

.summary-label {
    font-size: 12px;
    color: var(--sidebar-text);
}

.summary-value {
    font-size: 14px;
    font-weight: 700;
}

.summary-progress {
    background: rgba(255, 255, 255, 0.1);
}

.summary-detail {
    font-size: 11px;
    color: rgba(148, 163, 184, 0.6);
    margin-top: 6px;
}

.breadcrumb-current {
    color: var(--text-primary);
    font-weight: 500;
}

.topbar-left {
    display: flex;
    align-items: center;
    gap: 8px;
}

.burger-btn {
    display: none;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border: none;
    border-radius: 8px;
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    transition: background 0.15s;
}

.burger-btn:hover {
    background: #f1f5f9;
}

.burger-btn .material-icons-round {
    font-size: 22px;
}

.sidebar-overlay {
    display: none;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 99;
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
}

.sidebar-overlay.visible {
    opacity: 1;
    pointer-events: auto;
}

@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%);
    }

    .sidebar.open {
        transform: translateX(0);
    }

    .main {
        margin-left: 0;
    }

    .burger-btn {
        display: flex;
    }

    .sidebar-overlay {
        display: block;
    }

    .topbar {
        padding: 12px 16px;
    }

    .breadcrumb {
        display: none;
    }
}
</style>
