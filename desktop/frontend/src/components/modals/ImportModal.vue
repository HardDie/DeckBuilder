<template>
  <ui-modal
    v-model:show="isModalModel"
    title="Import game"
    @submit="onImport"
  >
    <div class="import-modal">
      <n-input
        v-model:value="form.name"
        placeholder="Title"
      />
      <n-upload
        ref="uploader"
        :on-change="onImportFile"
        :max="1"
        style="flex: 1"
        accept=".zip"
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
          <n-text style="font-size: 16px">Click or drag a file to this area to upload</n-text>
        </n-upload-dragger>
      </n-upload>
    </div>
  </ui-modal>
</template>

<script setup>
import { defineEmits, defineProps, reactive, ref, computed, watch } from 'vue'
import UiModal from '@/components/ui/uiModal.vue'
import { ArchiveOutlined } from '@vicons/material'

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

const form = reactive({
  file: null,
  name: '',
})

const uploader = ref(null)

const clearForm = () => {
  form.file = null
  form.name = ''
}

const onImportFile = ({ file }) => {
  form.file = file.file
}

watch(isModalModel, val => {
  if (!val) {
    clearForm()
    uploader.value.clear()
  }
})

const onImport = () => {
  if (form.file) {
    const formData = new FormData()
    formData.append('name', form.name)
    formData.append('file', form.file)
    emit('submit', formData)
  }
}
</script>
<style lang="scss">
.import-modal {
  padding: 5% 10%;
  height: 100%;
  display: flex;
  justify-content: center;
  flex-direction: column;
  gap: 30px;
  box-sizing: border-box;
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
