<template>
  <ui-modal
    v-model:show="isModalModel"
    title="Duplicate game"
    @submit="onDuplicate"
  >
    <div class="duplicate-modal">
      <img
        class="game-preview"
        :src="props.entityData.cachedImage"
        alt="preview"
      />
      <span class="duplicate-modal__text">{{ props.entityData.name }}</span>
      <span class="duplicate-modal__text">{{ props.entityData.description }}</span>
      <n-input
        v-model:value="saveName"
        placeholder="Save game as..."
      />
    </div>
  </ui-modal>
</template>

<script setup>
import UiModal from '@/components/ui/uiModal.vue'
import { computed, defineEmits, defineProps, ref } from 'vue'

const emit = defineEmits(['update:show', 'submit'])
const props = defineProps({
  show: {
    type: Boolean,
    default: false,
  },
  entityData: {
    type: Object,
    default: () => {},
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

const saveName = ref('')

const onDuplicate = () => {
  emit('submit', saveName.value)
}
</script>
<style lang="scss">
.duplicate-modal {
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  &__text {
    margin-bottom: 5px;
    font-weight: bold;
  }
}

.game-preview {
  max-width: 205px;
  border: 2px #138b44 solid;
  border-radius: 8px;
  user-select: none;
  &--unsetted {
    width: 205px;
    aspect-ratio: 0.71;
    background-color: #138b44;
  }
}
</style>
