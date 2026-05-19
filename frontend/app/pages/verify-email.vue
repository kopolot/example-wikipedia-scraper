<template>
    <div class="container py-5">
        <div class="row justify-content-center">
            <div class="col-md-6 col-lg-5">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h2 class="mb-4">{{ $t('verify.title') }}</h2>
                        <div v-if="loading" class="mb-3">
                            <div class="spinner-border text-primary" role="status">
                                <span class="visually-hidden">{{ $t('verify.loading') }}</span>
                            </div>
                        </div>
                        <div v-else-if="success" class="alert alert-success">
                            {{ $t('verify.success') }}<br>
                            <small>{{ $t('verify.redirect') }}</small>
                        </div>
                        <div v-else class="alert alert-danger">
                            {{ $t('verify.error') }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const route = useRoute();
const loading = ref(true);
const success = ref(false);
const token = route.query.token as string;

if (!token) {
    navigateTo('/');
}

onMounted(() => {

    useAuth().verifyEmail(token).then((result) => {
        loading.value = false;
        if (result.success) {
            success.value = true;
            setTimeout(() => {
                navigateTo('/');
            }, 3000);
        } else {
            success.value = false;
        }
    }).catch(() => {
        loading.value = false;
        success.value = false;
    });
})
</script>