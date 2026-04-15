<template>
    <div>
        <div class="tree-item" :class="{ active: isSelected }" :style="indent" @click="toggle">
            <Icon v-if="isFolder" name="chevron_right" class="chevron" :class="{ open: expanded }" />
            <span v-else class="spacer"></span>
            <Icon v-if="isFolder" :name="expanded ? 'folder_open' : 'folder'" class="folder-icon" />
            <Icon v-else name="description" class="file-icon" />
            <span class="node-name" :title="node.name">{{ node.name }}</span>
            <span class="cov-mini" :class="covCls">{{ covPct }}%</span>
        </div>
        <div v-if="isFolder && expanded" class="tree-children">
            <TreeNode v-for="child in filteredChildren" :key="child.name" :node="child" :depth="depth + 1"
                :selected-path="selectedPath" :reveal-path="revealPath" :expand-command="expandCommand"
                @select="onChildSelect" />
        </div>
    </div>
</template>

<script>
import { covClass } from '@/utils/coverage'
import Icon from '@/components/Icon.vue'

export default {
    name: 'TreeNode',
    components: { Icon },
    props: {
        node: { type: Object, required: true },
        depth: { type: Number, default: 0 },
        selectedPath: { type: String, default: '' },
        revealPath: { type: String, default: '' },
        expandCommand: { type: Object, default: null },
    },
    emits: ['select'],
    data() {
        return {
            expanded: this.depth < 1,
        }
    },
    created() {
        if (this.revealPath && this.node.type === 'folder' && this.node.path
            && this.revealPath.startsWith(this.node.path + '/')) {
            this.expanded = true
        }
    },
    watch: {
        revealPath(path) {
            if (path && this.isFolder && this.node.path && path.startsWith(this.node.path + '/')) {
                this.expanded = true
            }
        },
        expandCommand(cmd) {
            if (!cmd || !this.isFolder) return
            this.expanded = cmd.action === 'expand'
        },
    },
    computed: {
        isFolder() {
            return this.node.type === 'folder'
        },
        isSelected() {
            return this.selectedPath === this.node.path
        },
        covCls() {
            if (this.isFolder) return covClass(this.node.coverage)
            return covClass(this.node.file.coverage)
        },
        covPct() {
            if (this.isFolder) return this.node.coverage.toFixed(0)
            return this.node.file.coverage.toFixed(0)
        },
        indent() {
            return { paddingLeft: (this.depth * 13) + 'px' }
        },
        filteredChildren() {
            return this.node.children || []
        },
    },
    methods: {
        toggle() {
            if (this.isFolder) {
                this.expanded = !this.expanded
            }
            this.$emit('select', this.node)
        },
        onChildSelect(node) {
            this.$emit('select', node)
        },
    },
}
</script>

<style scoped>
.spacer {
    width: 18px;
}

.node-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tree-item {
    display: flex;
    align-items: center;
    gap: 2px;
    padding: 5px 8px;
    cursor: pointer;
    font-size: 13px;
    color: var(--text-secondary);
    border-radius: 4px;
    user-select: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tree-item:hover {
    background: #f1f5f9;
    color: var(--text-primary);
}

.tree-item.active {
    background: #e0f2fe;
    color: #0369a1;
}

.tree-item .icon {
    font-size: 16px;
    flex-shrink: 0;
}

.tree-item .chevron {
    font-size: 18px;
    color: var(--text-muted);
    width: 18px;
    transition: transform 0.15s;
}

.tree-item .chevron.open {
    transform: rotate(90deg);
}

.tree-item .folder-icon {
    color: #f59e0b;
}

.tree-item .file-icon {
    color: #3b82f6;
}

.tree-item .cov-mini {
    margin-left: auto;
    font-size: 11px;
    font-weight: 600;
    padding: 1px 6px;
    border-radius: 4px;
    flex-shrink: 0;
}

.tree-item .cov-mini.green {
    background: var(--green-bg);
    color: #065f46;
}

.tree-item .cov-mini.yellow {
    background: var(--yellow-bg);
    color: #92400e;
}

.tree-item .cov-mini.red {
    background: var(--red-bg);
    color: #991b1b;
}

.tree-children {
    padding-left: 0;
}
</style>
