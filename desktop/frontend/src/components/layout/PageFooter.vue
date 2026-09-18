<template>
  <div class="page-footer">
    <div class="page-footer__meta">
      <span>Total items: {{ itemsStore.itemsMeta.total }}</span>
      <span v-if="itemsStore.itemsMeta.cardsTotal">
        <span>, cards: {{ itemsStore.itemsMeta.cardsTotal }}</span>
      </span>
    </div>
    <span>{{ systemStore.version }}</span>
  </div>
</template>

<script setup>
import { useSystemStore } from '@/stores/system'
import { useItemsStore } from '@/stores/items'

const systemStore = useSystemStore()
const itemsStore = useItemsStore()
</script>

<style lang="scss">
.page-footer {
  background: #e3ded6;
  padding: 0 10px;
  display: flex;
  justify-content: space-between;
}
</style>
