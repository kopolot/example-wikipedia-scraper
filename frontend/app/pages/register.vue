
<template>
	<div class="container py-5">
		<div class="row justify-content-center">
			<div class="col-md-6 col-lg-5">
				<div class="card shadow-sm">
					<div class="card-body">
						<h2 class="mb-4 text-center">{{ $t('register.title') }}</h2>
						<form @submit.prevent="onSubmit">
							<div class="mb-3">
								<label for="username" class="form-label">{{ $t('register.username') }}</label>
								<input id="username" v-model="username" type="text" class="form-control" required>
							</div>
							<div class="mb-3">
								<label for="email" class="form-label">{{ $t('register.email') }}</label>
								<input id="email" v-model="email" type="email" class="form-control" required>
							</div>
							<div class="mb-3">
								<label for="password" class="form-label">{{ $t('register.password') }}</label>
								<input id="password" v-model="password" type="password" class="form-control" required>
							</div>
							<div class="mb-3">
								<label for="confirm" class="form-label">{{ $t('register.confirm') }}</label>
								<input id="confirm" v-model="confirm" type="password" class="form-control" required>
							</div>
							<div v-if="error" class="alert alert-danger text-center py-2">{{ error }}</div>
							<div v-if="success" class="alert alert-success text-center py-2">{{ $t('register.success') }}</div>
							<button type="submit" class="btn btn-primary w-100" :disabled="loading">
								{{ loading ? $t('register.submited') + '...' : $t('register.submit') }}
							</button>
						</form>
						<div class="text-center mt-3">
							<NuxtLink to="/">{{ $t('register.have_account') }}</NuxtLink>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
const loading = ref(false);
const success = ref(false);
const error = ref<string | null>(null);
const username = ref('');
const email = ref('');
const password = ref('');
const confirm = ref('');

const onSubmit = async () => {
	error.value = null;
	success.value = false;
	if (!username.value || username.value.length < 3) {
		error.value = $t('register.invalid_username');
		return;
	}
	if (!email.value.includes('@')) {
		error.value = $t('register.invalid_email');
		return;
	}
	if (password.value.length < 6) {
		error.value = $t('register.too_short');
		return;
	}
	if (password.value !== confirm.value) {
		error.value = $t('register.not_match');
		return;
	}
	loading.value = true;
	try {
		const result = await useAuth().register(username.value, email.value, password.value, confirm.value);
		loading.value = false;
		if (result.success) {
			success.value = true;
			setTimeout(() => {
				navigateTo('/');
			}, 3000);
		} else {
			error.value = result.errors ? result.errors.join(', ') : $t('register.error');
		}
	} catch {
		loading.value = false;
		error.value = $t('register.error');
	}
};
</script>

<style scoped>
/* Możesz dodać własne poprawki do stylu Bootstrapa */
</style>