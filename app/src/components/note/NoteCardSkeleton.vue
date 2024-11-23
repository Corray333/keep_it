<script lang="ts" setup>
import Skeleton from 'primevue/skeleton';




</script>

<template>
    
        <div class="note-card">
            <div class="note-card-header">
    
                <Skeleton height="1.5rem" class=" max-w-52 w-full"></Skeleton>
    
            </div>
    
            <div class="note-card-body">
                <Skeleton height="100%" class=" max-w-28 w-full"></Skeleton>
    
                <div class="note-card-skeleton-body-content">
                    <Skeleton height="1.5rem" class=" max-w-28 w-full"></Skeleton>
                    <Skeleton height="1.5rem" class=" max-w-28 w-full"></Skeleton>
                    <Skeleton height="1.5rem" class=" max-w-28 w-full"></Skeleton>
                </div>
            </div>
        </div>

</template>


<style scoped>

.note-card{
    @apply relative z-0 bg-primary-bg rounded-2xl shadow-lg border-2 border-invert-bg-50 w-full cursor-pointer;
    @apply duration-300 hover:scale-105
}

.note-card-header{
    @apply p-4 border-b-2 border-invert-bg-50 flex items-center justify-between;
}

.note-card-body{
    @apply p-4 flex h-36 overflow-hidden gap-2 text-sm;
}

.note-card-body-cover{
    @apply w-28 h-28 object-cover aspect-square rounded-lg;
}

.note-card-header-label{
    @apply flex gap-2 font-medium
}

.note-card-skeleton-body-content{
    @apply flex flex-col gap-2
}

</style>