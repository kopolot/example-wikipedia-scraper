
<template>
	<div class="container py-5">
		<div class="row justify-content-center">
			<div class="col-md-6 col-lg-5">
				<div class="card shadow-sm">
					<div class="card-body">
						<h2 class="mb-4 text-center">{{ $t('reset.title') }}</h2>
						<form @submit.prevent="onSubmit">
							<div class="mb-3">
								<label for="password" class="form-label">{{ $t('reset.password') }}</label>
								<input id="password" v-model="password" type="password" class="form-control" required>
							</div>
							<div class="mb-3">
								<label for="confirm" class="form-label">{{ $t('reset.confirm') }}</label>
								<input id="confirm" v-model="confirm" type="password" class="form-control" required>
							</div>
							<div v-if="error" class="alert alert-danger text-center py-2">{{ error }}</div>
							<div v-if="success" class="alert alert-success text-center py-2">{{ $t('reset.success') }}</div>
							<button type="submit" class="btn btn-primary w-100" :disabled="loading">
								{{ loading ? $t('reset.submited') + '...' : $t('reset.submit') }}
							</button>
						</form>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
const route = useRoute();
const loading = ref(false);
const success = ref(false);
const error = ref<string | null>(null);
const password = ref('');
const confirm = ref('');
const token = route.query.token as string;

if (!token) {
    navigateTo('/');
}

const onSubmit = async () => {
	error.value = null;
    success.value = false;
	if (password.value.length < 6) {
		error.value = $t('reset.too_short');
		return;
	}
	if (password.value !== confirm.value) {
		error.value = $t('reset.not_match');
		return;
	}
	loading.value = true;
	try {
		const result = await useAuth().resetPassword(token, password.value);
		loading.value = false;
		if (result.success) {
			success.value = true;
			setTimeout(() => {
				navigateTo('/');
			}, 3000);
		} else {
			error.value = result.errors?.join(', ') || $t('reset.error');
		}
	} catch {
		loading.value = false;
		error.value = $t('reset.error');
	}
};

</script>
