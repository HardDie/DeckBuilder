<template>
  <div
    class="main"
    oncontextmenu="return false;"
  >
    <add-modal
      v-model:show="isAddModal"
      :initial-data="tempItem"
      :entity-type="mainStore.itemType"
      :count="mainStore.itemType === 'cards'"
      :variables="mainStore.itemType === 'cards'"
      @submit="onAdd"
    />
    <import-modal
      v-model:show="isImportModal"
      @submit="onImport"
    />
    <duplicate-modal
      v-model:show="isDuplicateModal"
      :entity-data="tempItem"
      @submit="onDuplicate"
    />
    <replace-modal
      v-model:show="isReplaceModal"
      @submit="onReplace"
    />
    <page-header
      show-buttons
      :with-import="mainStore.itemType === 'games'"
      :with-replace="mainStore.itemType === 'games'"
      @on-add="openAddModal"
      @on-import="isImportModal = true"
      @on-filters="onFilters"
      @on-replace="isReplaceModal = true"
    />
    <page-content ref="contentRef">
      <transition-group name="slide-fade">
        <div
          v-if="items.length && !itemsStore.isItemsLoading"
          class="main__list"
        >
          <card-el
            v-for="item in items"
            :id="String(item.id)"
            :key="item.id"
            :name="item.name"
            :img="item.cachedImage"
            :description="item.description"
            :count="item.count"
            :clickable="mainStore.itemType !== 'cards'"
            :with-export="mainStore.itemType === 'games'"
            :with-duplicate="mainStore.itemType === 'games'"
            :with-render="mainStore.itemType === 'games'"
            :with-double-click="mainStore.itemType === 'cards'"
            @card-click="onItemClick"
            @on-export="onExport"
            @on-duplicate="openDuplicateModal"
            @on-render="onGenerate"
            @on-edit="openAddModal"
            @on-delete="onDeleteItem"
            @on-double-click="openAddModal"
          />
        </div>
      </transition-group>
      <div
        v-if="!items.length && !isItemsLoading"
        class="empty-filler"
      >
        There is no items
      </div>
    </page-content>
    <page-footer />
    <div
      v-if="generateProgress"
      class="render-spinner"
    >
      <n-progress
        class="render-spinner__progress"
        type="circle"
        color="#169747"
        :percentage="generateProgress"
      />
    </div>
    <div
      v-if="itemsStore.isApiPending"
      class="render-spinner"
    >
      <n-spin color="#169747" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, onBeforeUnmount, h } from 'vue'
import { useItemsStore } from '@/stores/items'
import { useStore } from '@/stores/main'
import { useSystemStore } from '@/stores/system'
import { useReplaceStore } from '@/stores/replace'
import { useRouter, useRoute } from 'vue-router'

import PageHeader from '@/components/layout/PageHeader.vue'
import PageContent from '@/components/layout/PageContent.vue'
import AddModal from '@/components/modals/AddModal.vue'
import DuplicateModal from '@/components/modals/DuplicateModal.vue'
import ImportModal from '@/components/modals/ImportModal.vue'
import ReplaceModal from '@/components/modals/ReplaceModal.vue'
import CardEl from '@/components/CardEl.vue'
import { useDialog } from 'naive-ui'
import PageFooter from '@/components/layout/PageFooter.vue'

const router = useRouter()
const route = useRoute()

const itemsStore = useItemsStore()
const mainStore = useStore()
const systemStore = useSystemStore()
const replaceStore = useReplaceStore()

const isAddModal = ref(false)
const isImportModal = ref(false)
const isDuplicateModal = ref(false)
const isReplaceModal = ref(false)

const items = computed(() => itemsStore.items)
const isItemsLoading = computed(() => itemsStore.isItemsLoading)

const tempItem = ref({})

const contentRef = ref(null)

onMounted(() => {
  mainStore.setItemType(route.name)
  itemsStore.fetchItems({ ...route.params })
  startCheckStatusPolling()
})

watch(isAddModal, val => {
  if (!val) {
    tempItem.value = {}
  }
})

watch(isDuplicateModal, val => {
  if (!val) {
    tempItem.value = {}
  }
})

const currentPathId = computed(() => {
  let thisId

  switch (mainStore.itemType) {
    case 'games':
      thisId = 'gameId'
      break
    case 'collections':
      thisId = 'collectionId'
      break
    case 'decks':
      thisId = 'deckId'
      break
    case 'cards':
      thisId = 'cardId'
      break
    default:
      thisId = undefined
  }
  return thisId
})

const openAddModal = id => {
  if (id) {
    itemsStore
      .fetchItem({
        ...route.params,
        [currentPathId.value]: id,
      })
      .then(response => {
        tempItem.value = response
        isAddModal.value = true
      })
  } else {
    isAddModal.value = true
  }
}

const onAdd = ({ mode, data }) => {
  const action = mode === 'create' ? itemsStore.fetchCreateItem : itemsStore.fetchEditItem
  action({
    body: data,
    ...route.params,
    [currentPathId.value]: tempItem.value?.id,
    id: tempItem.value?.id,
  })
    .then(() => {
      if (mode === 'create') {
        contentRef.value.$el.scrollTop = contentRef.value.$el.scrollHeight
      }
    })
    .finally(() => {
      isAddModal.value = false
    })
}

const openDuplicateModal = id => {
  itemsStore.fetchItem({ [currentPathId.value]: id }).then(response => {
    tempItem.value = response
    isDuplicateModal.value = true
  })
}

const onDuplicate = name => {
  itemsStore
    .fetchDuplicateGame({
      [currentPathId.value]: tempItem.value.id,
      body: { name },
    })
    .finally(() => {
      isDuplicateModal.value = false
    })
}

const dialog = useDialog()
const onDeleteItem = id => {
  const deletedName = items.value?.find(el => String(el.id) === id).name
  dialog.error({
    title: 'Delete confirmation',
    content: () =>
      h('span', {}, [
        `Are you sure you want to permanently delete `,
        h('span', { style: { fontWeight: 'bold' } }, deletedName),
        ' ?',
      ]),
    positiveText: 'Delete',
    negativeText: 'Cancel',
    showIcon: false,
    onPositiveClick: () =>
      itemsStore.fetchDeleteItem({
        ...route.params,
        [currentPathId.value]: id,
        id,
      }),
  })
}

const onImport = data => {
  itemsStore.fetchImportGame(data).then(() => {
    isImportModal.value = false
  })
}

const onExport = id => {
  itemsStore.fetchExportGame({ [currentPathId.value]: id })
}

const generateInterval = ref(null)
const generateProgress = ref(0)

const startCheckStatusPolling = () => {
  generateInterval.value = setInterval(() => {
    systemStore.fetchCheckStatus().then(status => {
      generateProgress.value = Math.floor(status.progress)
      if (status.status !== 'in_progress') {
        clearInterval(generateInterval.value)
        setTimeout(() => {
          generateProgress.value = 0
        }, 500)
      }
    })
  }, 500)
}

const onGenerate = id => {
  itemsStore.fetchGenerateGame({ [currentPathId.value]: id }).then(() => {
    startCheckStatusPolling()
  })
}

const onFilters = () => {
  itemsStore.fetchItems({ ...route.params })
}

const onReplace = data => {
  const formData = new FormData()
  formData.append('file', data.file)
  const jsonBlob = new Blob([JSON.stringify({ data: data.mapping })], { type: 'application/json' })
  formData.append('mapping', jsonBlob)
  replaceStore.replace(formData)
}

const onItemClick = id => {
  router.push(`${route.fullPath}/${mainStore.itemType}/${id}`.replace('//', '/'))
}

onBeforeUnmount(() => {
  itemsStore.clearItems()
})
</script>

<style lang="scss">
.main {
  height: 100%;
  display: flex;
  flex-direction: column;

  &__list {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: center;
  }
}

.empty-filler {
  position: absolute;
  top: calc(50% + 28px);
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 2rem;
  color: grey;
  user-select: none;
}

.render-spinner {
  z-index: 999;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(227, 222, 214, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
