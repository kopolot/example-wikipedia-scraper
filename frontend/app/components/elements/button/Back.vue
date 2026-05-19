<template>
    <div class="w-100 d-flex align-items-center">
        <button class="btn btn-link p-0 text-primary d-flex align-items-center gap-1" @click="$router.back()">
            <BootstrapIcon name="arrow-left"/>
            Powrót
        </button>
    </div>
</template>

<style scoped lang="scss">
button {
    margin-bottom: 1rem;
}
</style>