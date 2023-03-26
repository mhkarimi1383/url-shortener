<template id="app">
  <div>
    <NuxtPage />
  </div>
</template>
<style lang="scss">
* {

  &::-webkit-scrollbar {
    width: 4px;
    height: 8px;
  }

  &::-webkit-scrollbar-track {
    box-shadow: none;
    border-radius: 10px;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--va-primary);
    border-radius: 2px;
  }

  scrollbar-color: var(--va-primary) transparent;
  scrollbar-width: thin;
}
</style>
