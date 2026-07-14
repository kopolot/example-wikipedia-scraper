<template>
    <div class="container-fluid py-4">
        <div class="row justify-content-center">
            <div class="col-12">
                <div class="mb-4 text-center">
                    <h1 class="mb-2">{{ $t('panel.user.title') }}</h1>
                    <h3 class="text-muted">{{ $t('panel.user.greetings', { name: user?.username }) }}</h3>
                </div>
                <div class="row g-4">
                    <div class="col-md-4">
                        <div class="card shadow-sm mb-4">
                            <div class="card-body">
                                <h2 class="h5 mb-3">{{ $t('panel.user.dataTitle') }}</h2>
                                <p class="mb-2"><strong>{{ $t('panel.user.username') }}:</strong> {{ user?.username }}</p>
                                <p class="mb-2"><strong>{{ $t('panel.user.email') }}:</strong> {{ user?.email }}</p>
                                <p class="mb-0"><strong>{{ $t('panel.user.registeredAt') }}:</strong> {{ user?.createdAt?.toLocaleString() }}</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card shadow-sm mb-4">
                            <div class="card-body text-center">
                                <h2 class="h5 mb-3">{{ $t('panel.user.subscriptionTitle') }}</h2>
                                <div v-if="user?.subscriptionLevel ?? 0> 0">
                                    <p>{{ $t('panel.user.subscriptionLevel') }}: <strong>{{ subscriptionInfo?.name }}</strong></p>
                                    <p v-if="user?.subscriptionExpiration" :class="{ 'text-danger': subExpired }">{{ subExpired ? $t('panel.user.subscriptionExpired') : $t('panel.user.subscriptionExpires') }}: <strong>{{ user?.subscriptionExpiration?.toLocaleString() }}</strong></p>
                                    <div class="d-flex justify-content-center gap-3 mt-4">
                                        <button v-if="user?.subscriptionLevel === 0" class="btn btn-success" @click="goToSubscribe">{{$t('panel.user.start')}}</button>
                                        <button v-else-if="((user?.subscriptionLevel ?? 0) > 0) && subExpired" class="btn btn-warning" @click="goToSubscribe">{{$t('panel.user.renewSubscription')}}</button>
                                        <!-- <button v-else class="btn btn-info" @click="goToSubscribe">{{$t('panel.user.renewSubscription')}}</button> -->
                                    </div>
                                </div>
                                <div v-else>
                                    <p>{{ $t('panel.user.noActiveSubscription') }}</p>
                                    <button class="btn btn-success" @click="goToSubscribe">{{$t('panel.user.startSubscription') }}</button>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4">
                        <div class="card shadow-sm mb-4">
                            <div class="card-body d-flex flex-column align-items-start gap-2">
                                <h2 class="h5 mb-3">{{ $t('panel.user.actionsTitle') }}</h2>
                                <NuxtLink class="btn btn-outline-primary w-100" to="user/change_password">
                                    {{ $t('panel.user.changePassword') }}
                                </NuxtLink>
                                <NuxtLink class="btn btn-outline-secondary w-100" to="user/change_email">
                                    {{ $t('panel.user.changeEmail') }}
                                </NuxtLink>
                                <button class="btn btn-outline-danger w-100" @click="handleLogoutAll">
                                    {{ $t('panel.user.logoutAll') }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
definePageMeta({
    middleware: 'auth',
    layout: 'panel'
})

const auth = useAuth();
const userComposable = useUser();
const user = auth.user;
const subscriptionInfo = userComposable.subscriptionData;
const router = useRouter();

const subExpired = computed(() => {
    if (!user.value || !user.value.subscriptionExpiration) return true;
    return new Date(user.value.subscriptionExpiration) < new Date();
});

onMounted(async () => {
    if (user.value?.subscriptionLevel === undefined) {
        await auth.loadUserData();
    }
    if ( subscriptionInfo.value === null) {
        userComposable.getSubscriptionData();
    }
}); 


async function handleLogoutAll() {
    const result = await userComposable.logoutAll();
    if (result.success) {
        await auth.logout();
    } else {
        // Możesz dodać obsługę błędu, np. toast
        console.error('Logout all failed', result.error);
    }
}

function goToSubscribe() {
    router.push({ path: '/panel/user/subscribe' });
}
</script>

<style lang="scss" scoped>
/* Możesz dodać własne style, jeśli chcesz nadpisać domyślne style Bootstrap */
</style>