<template>
    <div class="loading-screen">
        <div class="logo-big">Go</div>
        <h2>Go Coverage Dashboard</h2>
        <p>Load a Go <code>coverprofile.out</code> file to visualize code coverage with interactive dashboards and
            file-level details.</p>
        <div class="config-form">
            <div class="form-group">
                <label>Coverage Profile URL</label>
                <input v-model="profileUrl" placeholder="./cover.out" @keyup.enter="load">
                <div class="hint">Relative or absolute URL to the coverprofile file</div>
            </div>
            <div class="form-group">
                <label>Source Root URL (optional)</label>
                <input v-model="sourceRoot" placeholder="http://localhost:8080/">
                <div class="hint">Base URL to fetch Go source files for line-by-line view</div>
            </div>
            <div class="form-group">
                <label>Module Prefix (auto-detected)</label>
                <input v-model="modulePrefix" placeholder="github.com/user/project">
                <div class="hint">Will be stripped from file paths when fetching source</div>
            </div>
            <div class="form-group">
                <label>Route Coverage URL (optional)</label>
                <input v-model="routeCoverageUrl" placeholder="./route-coverage.rcov">
                <div class="hint">URL to a route-coverage.rcov file for API route coverage</div>
            </div>
            <div class="load-actions">
                <button class="btn btn-primary" @click="load" :disabled="loading">
                    <span v-if="loading" class="spinner spinner-sm"></span>
                    <span class="material-icons-round" v-else>play_arrow</span>
                    {{ loading ? 'Loading...' : 'Load Coverage' }}
                </button>
            </div>
            <div v-if="loadError" class="error-msg">{{ loadError }}</div>
        </div>
    </div>
</template>

<script>
import { parseCoverProfile, buildCoverData } from '@/utils/coverage'

export default {
    name: 'ConfigScreen',
    props: {
        initialProfileUrl: { type: String, default: './cover.out' },
        initialSourceRoot: { type: String, default: '' },
        initialModulePrefix: { type: String, default: '' },
        initialRouteCoverageUrl: { type: String, default: '' },
    },
    emits: ['loaded'],
    data() {
        return {
            profileUrl: this.initialProfileUrl,
            sourceRoot: this.initialSourceRoot,
            modulePrefix: this.initialModulePrefix,
            routeCoverageUrl: this.initialRouteCoverageUrl,
            loading: false,
            loadError: '',
        }
    },
    methods: {
        async load() {
            this.loading = true
            this.loadError = ''
            try {
                const resp = await fetch(this.profileUrl)
                if (!resp.ok) throw new Error('HTTP ' + resp.status + ': ' + resp.statusText)
                const text = await resp.text()
                const parsed = parseCoverProfile(text)
                const coverData = buildCoverData(parsed)
                if (!this.modulePrefix) this.modulePrefix = coverData.modulePrefix
                this.$emit('loaded', {
                    coverData,
                    profileUrl: this.profileUrl,
                    sourceRoot: this.sourceRoot,
                    modulePrefix: this.modulePrefix,
                    routeCoverageUrl: this.routeCoverageUrl,
                })
            } catch (e) {
                this.loadError = 'Failed to load coverage profile: ' + e.message
            } finally {
                this.loading = false
            }
        },
    },
}
</script>

<style scoped>
.loading-screen {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    width: 100%;
    text-align: center;
    padding: 40px;
}

.loading-screen .logo-big {
    width: 64px;
    height: 64px;
    background: var(--teal);
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 700;
    font-size: 24px;
    margin-bottom: 24px;
}

.loading-screen h2 {
    font-size: 24px;
    font-weight: 700;
    color: var(--text-primary);
    margin-bottom: 8px;
}

.loading-screen p {
    color: var(--text-secondary);
    font-size: 14px;
    margin-bottom: 28px;
    max-width: 420px;
}

.config-form {
    background: var(--card-bg);
    border-radius: var(--radius);
    box-shadow: var(--shadow-lg);
    padding: 28px;
    width: 100%;
    max-width: 480px;
    text-align: left;
    border: 1px solid var(--border);
}

.form-group {
    margin-bottom: 16px;
}

.form-group label {
    display: block;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
    margin-bottom: 6px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.form-group input {
    width: 100%;
    padding: 10px 14px;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 14px;
    font-family: 'JetBrains Mono', monospace;
    outline: none;
    background: #f8fafc;
}

.form-group input:focus {
    border-color: var(--teal);
    box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.12);
}

.form-group .hint {
    font-size: 11px;
    color: var(--text-muted);
    margin-top: 4px;
}

.load-actions {
    display: flex;
    gap: 10px;
    align-items: center;
}

.error-msg {
    background: var(--red-bg);
    color: #991b1b;
    padding: 10px 14px;
    border-radius: 8px;
    font-size: 13px;
    margin-top: 12px;
    border: 1px solid var(--red-border);
}

.spinner {
    width: 24px;
    height: 24px;
    border: 3px solid var(--border);
    border-top-color: var(--teal);
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
}

.spinner-sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
</style>
