<template>
  <ui-modal
    v-model:show="isModalModel"
    title="Replace"
    @submit="onReplace"
  >
    <div class="replace-modal">
      <div
        v-show="!currentFile"
        class="input-upload"
      >
        <n-upload
          ref="uploader"
          :on-change="onImportFile"
          :max="1"
          accept=".json"
        >
          <n-upload-dragger class="upload-dragger">
            <div style="margin-bottom: 12px">
              <n-icon
                size="48"
                :depth="3"
              >
                <archive-outlined />
              </n-icon>
            </div>
            <n-text style="font-size: 16px">
              Please upload your json file where you want to replace the images
            </n-text>
          </n-upload-dragger>
        </n-upload>
      </div>
      <div class="inputs">
        <div
          v-for="el in data"
          :key="el.key"
          class="inputs__element"
        >
          <span class="label">{{ el.key.split('/').at(-1) }}</span>
          <n-input
            v-model:value="el.value"
            placeholder="Image URL"
          />
        </div>
      </div>
    </div>
  </ui-modal>
</template>
<script setup>
import UiModal from '@/components/ui/uiModal.vue'
import { computed, defineEmits, defineProps, ref, watch } from 'vue'
import { ArchiveOutlined } from '@vicons/material'
import { useReplaceStore } from '@/stores/replace'

const replaceStore = useReplaceStore()

const emit = defineEmits(['update:show', 'submit'])

const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
})

const isModalModel = computed({
  get() {
    return props.show
  },
  set(val) {
    emit('update:show', val)
  },
})

const uploader = ref(null)
const currentFile = ref(null)
const data = ref([])

const onImportFile = ({ file }) => {
  currentFile.value = file.file
  const formData = new FormData()
  formData.append('file', currentFile.value)
  replaceStore
    .prepare(formData)
    .then(res => (data.value = res.data))
    .catch(() => {
      clearInput()
    })
}

const clearInput = () => {
  uploader.value.clear()
  currentFile.value = null
}

watch(isModalModel, val => {
  if (!val) {
    clearInput()
    data.value = []
  }
})

const onReplace = () => {
  emit('submit', { file: currentFile.value, mapping: data.value })
}
</script>

<style lang="scss">
.replace-modal {
  padding: 5% 10%;
  height: 100%;
  box-sizing: border-box;
}

.inputs {
  display: flex;
  flex-direction: column;
  gap: 10px;
  &__element {
    .label {
      font-size: 0.8rem;
      font-weight: normal;
    }
  }
}

.upload-dragger {
  background: #e3ded6;
  border-color: #ababab;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
</style>
