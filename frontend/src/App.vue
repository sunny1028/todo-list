<script setup lang="ts">
import AppHeader from './components/layout/AppHeader.vue'
import AppNav from './components/layout/AppNav.vue'
import Sidebar from './components/layout/Sidebar.vue'
import ToastContainer from './components/ui/ToastContainer.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 transition-colors">
    <AppHeader />
    <div class="flex">
      <Sidebar />
      <main class="flex-1 mx-auto max-w-2xl px-4 py-6 pb-20 md:pb-6">
        <RouterView />
      </main>
    </div>
    <AppNav />
    <ToastContainer />
  </div>
</template>
