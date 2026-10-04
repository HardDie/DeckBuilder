<template>
  <ui-modal
    v-model:show="isModalModel"
    title="Settings"
    @submit="onSave"
  >
    <n-form
      class="settings-modal"
      label-placement="left"
      label-width="auto"
      :disabled="!isLoaded"
    >
      <n-form-item label="Card scale:">
        <n-input-number
          v-model:value="form.card_scale"
          :min="minScale"
          :max="maxScale"
          :step="0.1"
          :precision="2"
        />
      </n-form-item>
      <div class="settings-modal__hint">
        Size of rendered cards and decks in Tabletop Simulator. 1 is the default size.
      </div>
      <n-form-item label="Back shadow:">
        <n-switch v-model:value="form.enable_back_shadow" />
      </n-form-item>
      <div class="settings-modal__hint">Darkens the card back on the sheet.</div>
      <n-form-item label="Log level:">
        <n-select
          v-model:value="form.log_level"
          :options="logLevels"
        />
      </n-form-item>
      <div class="settings-modal__hint">
        Debug also logs how long each render step takes. Warn and Error write only problems. Applies
        at once.
      </div>
    </n-form>
  </ui-modal>
</template>

<script setup>
import UiModal from '@/components/ui/uiModal.vue'
import { computed, reactive, ref, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useSystemStore } from '@/stores/system'

// Same limits as entities/settings (MinCardScale, MaxCardScale).
const minScale = 0.1
const maxScale = 10
// Same values as entities/settings (LogLevel*), lowest first.
const logLevels = [
  { label: 'Debug', value: 'debug' },
  { label: 'Info', value: 'info' },
  { label: 'Warn', value: 'warn' },
  { label: 'Error', value: 'error' },
]

const emit = defineEmits(['update:show'])
const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
})

const systemStore = useSystemStore()
const toast = useToast()

const isModalModel = computed({
  get() {
    return props.show
  },
  set(val) {
    emit('update:show', val)
  },
})

const isLoaded = ref(false)
const form = reactive({ card_scale: 1, enable_back_shadow: false, log_level: 'info' })

watch(
  () => props.show,
  val => {
    if (!val) {
      return
    }
    isLoaded.value = false
    systemStore
      .fetchSettings()
      .then(settings => {
        form.card_scale = settings.card_scale
        form.enable_back_shadow = settings.enable_back_shadow
        form.log_level = settings.log_level
        isLoaded.value = true
      })
      .catch(() => {
        isModalModel.value = false
      })
  },
)

const onSave = () => {
  if (!isLoaded.value) {
    return
  }
  const scale = form.card_scale
  if (typeof scale !== 'number' || scale < minScale || scale > maxScale) {
    toast.error(`Card scale must be between ${minScale} and ${maxScale}`)
    return
  }
  systemStore
    .saveSettings({
      card_scale: scale,
      enable_back_shadow: form.enable_back_shadow,
      log_level: form.log_level,
    })
    .then(() => {
      toast.success('Settings saved')
      isModalModel.value = false
    })
    .catch(() => {})
}
</script>

<style lang="scss">
.settings-modal {
  padding: 20px;

  &__hint {
    margin: -16px 0 16px;
    font-size: 12px;
    color: #5a5348;
  }
}
</style>
