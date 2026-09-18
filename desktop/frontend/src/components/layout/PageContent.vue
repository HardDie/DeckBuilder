<template>
  <div class="page-content">
    <slot />
  </div>
</template>

<style lang="scss">
.page-content {
  background-color: #e3ded6;
  flex: 1 1 0;
  padding: 15px;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
