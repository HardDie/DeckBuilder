<template>
  <div :class="['card', { 'card--hover': props.clickable }]">
    <n-tooltip
      trigger="hover"
      placement="bottom"
      :delay="1000"
      :disabled="!propsRef.description.value"
    >
      <template #trigger>
        <div class="img-wrap">
          <img
            ref="cardImg"
            class="img"
            rel="preload"
            :src="propsRef.missingImage.value ? noImageSrc : propsRef.img.value"
            :alt="propsRef.name.value"
            @contextmenu="onRightClick"
            @click="onCardClick"
          />
          <n-tooltip
            v-if="propsRef.warning.value"
            trigger="hover"
            placement="top"
          >
            <template #trigger>
              <n-icon
                class="img-wrap__warning"
                size="28"
                color="#e0a000"
              >
                <warning-amber-round />
              </n-icon>
            </template>
            <span>{{ propsRef.warning.value }}</span>
          </n-tooltip>
        </div>
      </template>
      <span>{{ propsRef.description.value }}</span>
    </n-tooltip>
    <div class="label">
      <span class="label__name">{{ propsRef.name.value }}</span>
      <span class="label__count">{{ cardCount }}</span>
    </div>
    <n-dropdown
      placement="bottom-start"
      trigger="manual"
      :x="xRef"
      :y="yRef"
      :options="options"
      :show="showDropdownRef"
      :on-clickoutside="onClickOutside"
      @select="handleSelect"
    />
  </div>
</template>

<script setup>
import { ref, computed, toRefs, onMounted, onBeforeUnmount } from 'vue'
import { WarningAmberRound } from '@vicons/material'

// Shown instead of a broken image when a deck or card has no image yet.
const noImageSrc =
  'data:image/svg+xml;utf8,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="205" height="288" viewBox="0 0 205 288">' +
      '<rect width="205" height="288" fill="#e3ded6"/>' +
      '<text x="102" y="150" font-family="Roboto, sans-serif" font-size="18" fill="#7a7a7a" text-anchor="middle">No image</text>' +
      '</svg>',
  )

const emit = defineEmits([
  'card-click',
  'on-edit',
  'on-delete',
  'on-export',
  'on-render',
  'on-duplicate',
  'on-double-click',
])

const props = defineProps({
  name: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
  img: {
    type: String,
    default: '',
  },
  id: {
    type: String,
    required: true,
  },
  withExport: {
    type: Boolean,
    default: false,
  },
  withDuplicate: {
    type: Boolean,
    default: false,
  },
  withRender: {
    type: Boolean,
    default: false,
  },
  withDoubleClick: {
    type: Boolean,
    default: false,
  },
  count: {
    type: Number,
    default: 1,
  },
  clickable: {
    type: Boolean,
    default: false,
  },
  // The item's own image is missing: show a placeholder instead of a broken image.
  missingImage: {
    type: Boolean,
    default: false,
  },
  // Shown as a warning badge; empty means no badge.
  warning: {
    type: String,
    default: '',
  },
})

const onDoubleClick = () => {
  emit('on-double-click', propsRef.id.value)
}

onMounted(() => {
  if (!propsRef.withDoubleClick.value) {
    return
  }
  cardImg.value.addEventListener('dblclick', onDoubleClick)
})

onBeforeUnmount(() => {
  if (!propsRef.withDoubleClick.value) {
    return
  }
  cardImg.value.removeEventListener('dblclick', onDoubleClick)
})

const cardImg = ref(null)

const propsRef = toRefs(props)

const cardCount = computed(() => (propsRef.count.value > 1 ? ` x${propsRef.count.value}` : ''))

// context menu
const showDropdownRef = ref(false)
const xRef = ref(0)
const yRef = ref(0)

const options = computed(() => {
  const optionsArr = [
    {
      label: 'Change',
      key: 'change',
    },
  ]
  if (propsRef.withExport.value) {
    optionsArr.push({
      label: 'Export',
      key: 'export',
    })
  }
  if (propsRef.withDuplicate.value) {
    optionsArr.push({
      label: 'Duplicate',
      key: 'duplicate',
    })
  }
  if (propsRef.withRender.value) {
    optionsArr.push({
      label: 'Render',
      key: 'render',
    })
  }
  optionsArr.push({
    label: 'Delete',
    key: 'delete',
  })
  return optionsArr
})

const onRightClick = e => {
  e.preventDefault()
  showDropdownRef.value = true
  xRef.value = e.clientX
  yRef.value = e.clientY
}
const onClickOutside = () => {
  showDropdownRef.value = false
}

const handleSelect = key => {
  if (key === 'change') {
    emit('on-edit', propsRef.id.value)
  }
  if (key === 'delete') {
    emit('on-delete', propsRef.id.value)
  }
  if (key === 'export') {
    emit('on-export', propsRef.id.value)
  }
  if (key === 'duplicate') {
    emit('on-duplicate', propsRef.id.value)
  }
  if (key === 'render') {
    emit('on-render', propsRef.id.value)
  }
  showDropdownRef.value = false
}

const onCardClick = () => {
  if (props.clickable) {
    emit('card-click', propsRef.id.value)
  }
}
</script>

<style lang="scss">
.img-wrap {
  position: relative;

  &__warning {
    position: absolute;
    top: 8px;
    right: 8px;
    filter: drop-shadow(0 0 2px rgba(0, 0, 0, 0.5));
  }
}

.img {
  max-width: 205px;
  border: 2px #138b44 solid;
  border-radius: 8px;
  user-select: none;
  cursor: pointer;
}

.label {
  font-family: 'Roboto', sans-serif;
  font-size: 16px;
  color: black;
  font-weight: bold;
  text-align: center;
  margin-top: 8px;

  &__count {
    color: grey;
  }
}

.card {
  width: min-content;

  &--hover:hover {
    transform: scale(1.03);
  }

  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;
}
</style>
