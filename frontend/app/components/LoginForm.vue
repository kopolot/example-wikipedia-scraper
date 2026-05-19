<template>
    <div class="container py-5">
        <div class="row justify-content-center">
            <div class="col-md-6 col-lg-5">
                <div class="card shadow-sm">
                    <div class="card-body">
                        <form @submit.prevent="onLogin">
                            <h2 class="mb-4 text-center">{{ $t('login.title') }}</h2>
                            <div class="mb-3">
                                <label for="email" class="form-label">{{ $t('login.email') }}</label>
                                <input id="email" v-model="email" type="email" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">{{ $t('login.password') }}</label>
                                <input id="password" v-model="password" type="password" class="form-control" required>
                            </div>
                            <div v-if="error" class="alert alert-danger text-center py-2">{{ error }}</div>
                            <div class="py-2 text-center">
                                <NuxtLink to="/forgot-password">{{ $t('login.forgot_password') }}</NuxtLink>
                            </div>
              
                            <button type="submit" class="btn btn-primary w-100" :disabled="loading">
                                {{ loading ? $t('login.submited') + '...' : $t('login.submit') }}
                            </button>
                                          <div class="text-center mt-2">
                                <NuxtLink to="/register">{{ $t('login.no_account') }}</NuxtLink>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const auth = useAuth(),
    t = useI18n().t,
    loading = ref(false),
    email = ref(''),
    password = ref(''),
    error = ref<string | null>(null),
    onLogin = async () => {
        try {
            loading.value = true;
            const {success, errors} = await auth.login(email.value, password.value);
            loading.value = false;
            if (success) {
                navigateTo('/panel');
            } else {
                error.value = errors ? errors.join(', ') : t('login.error');
            }
        } catch (e) {
            console.error('Login exception:', e);
            loading.value = false;
            error.value = t('login.error');
        }
    };
</script>
