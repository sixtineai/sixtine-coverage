<template>
    <div class="explorer-layout">
        <div class="tree-panel">
            <div class="tree-search">
                <input v-model="treeSearch" placeholder="Filter files...">
            </div>
            <div class="tree-nodes">
                <TreeNode v-for="child in filteredTree" :key="child.name" :node="child" :depth="0"
                    :selected-path="selectedPath" @select="onTreeSelect" />
            </div>
        </div>
        <div class="tree-content">
            <!-- File Coverage View -->
            <FileCoverage v-if="currentView === 'filecov' && selectedFile" :file="selectedFile"
                :source-root="sourceRoot" :module-prefix="modulePrefix" />
            <!-- Folder View -->
            <FolderView v-else :folder="selectedFolder" @open-file="onOpenFile" />
        </div>
    </div>
</template>

<script>
import TreeNode from '@/components/TreeNode.vue'
import FileCoverage from '@/components/FileCoverage.vue'
import FolderView from '@/components/FolderView.vue'
import { filterTreeNodes, findTreeNode } from '@/utils/coverage'

export default {
    name: 'TreeExplorer',
    components: { TreeNode, FileCoverage, FolderView },
    props: {
        tree: { type: Array, required: true },
        sourceRoot: { type: String, default: '' },
        modulePrefix: { type: String, default: '' },
        initialFolderPath: { type: String, default: '' },
    },
    emits: ['view-change'],
    data() {
        return {
            treeSearch: '',
            selectedPath: '',
            selectedFile: null,
            selectedFolder: null,
            currentView: 'tree', // 'tree' | 'folder' | 'filecov'
        }
    },
    computed: {
        filteredTree() {
            if (!this.treeSearch) return this.tree
            const q = this.treeSearch.toLowerCase()
            return filterTreeNodes(this.tree, q)
        },
    },
    watch: {
        initialFolderPath: {
            handler(path) {
                if (path) {
                    const node = findTreeNode(this.tree, path)
                    if (node) this.onTreeSelect(node)
                }
            },
            immediate: true,
        },
    },
    methods: {
        onTreeSelect(node) {
            this.selectedPath = node.path
            if (node.type === 'folder') {
                this.selectedFolder = node
                this.selectedFile = null
                this.currentView = 'folder'
                this.$emit('view-change', { view: 'folder', node })
            } else if (node.type === 'file') {
                this.onOpenFile(node.file)
            }
        },
        onOpenFile(file) {
            this.selectedFile = file
            this.selectedPath = file.relativePath || file.path
            this.currentView = 'filecov'
            this.$emit('view-change', { view: 'filecov', file })
        },
    },
}
</script>

<style scoped>
.explorer-layout {
    display: flex;
    height: calc(100vh - 49px);
}

.tree-nodes {
    padding: 8px;
}

.tree-panel {
    width: 320px;
    min-width: 280px;
    border-right: 1px solid var(--border);
    background: var(--card-bg);
    overflow-y: auto;
    max-height: calc(100vh - 50px);
    position: sticky;
    top: 49px;
}

.tree-content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

.tree-search {
    padding: 10px;
    border-bottom: 1px solid var(--border);
}

.tree-search input {
    width: 100%;
    padding: 7px 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    font-size: 12px;
    font-family: inherit;
    outline: none;
    background: #f8fafc;
}

.tree-search input:focus {
    border-color: var(--teal);
    box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.1);
}

@media (max-width: 768px) {
    .explorer-layout {
        flex-direction: column;
        height: auto;
    }

    .tree-panel {
        width: 100%;
        min-width: 0;
        max-height: 260px;
        position: static;
        border-right: none;
        border-bottom: 1px solid var(--border);
    }

    .tree-content {
        padding: 16px 14px;
    }
}
</style>
