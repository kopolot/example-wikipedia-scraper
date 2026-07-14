<template>
  <div class="container py-4">
    <div class="text-center mb-3">
      <NuxtLink to="/panel/user" class="btn btn-outline-secondary">
        {{$t('back_to_panel') || 'Powrót do panelu'}}
      </NuxtLink>
    </div>
    <div class="row justify-content-center">
      <div class="col-md-8">
        <div class="card shadow-sm mb-4">
          <div class="card-body">
            <h2 class="h4 mb-4 text-center">{{ $t('choose_subscription_level') }}</h2>
            <div v-for="([levelId, products]) in Array.from(subscriptionLevels)" :key="levelId" class="mb-4">
              <div class="mb-2 fw-bold">
                {{ products && products.length > 0 && products[0]?.subscriptionLevel ? products[0].subscriptionLevel.name : 'Poziom ' + levelId }}
                <span v-if="products && products.length > 0 && products[0]?.subscriptionLevel" class="ms-2 text-muted small">(Limit: {{ products[0].subscriptionLevel.limit }})</span>
              </div>
              <div class="row g-3">
                <div v-for="prod in products" :key="prod.productId" class="col-md-6">
                  <div class="card h-100" :class="{ 'border-primary': selectedLevel === levelId && selectedProduct?.id === prod.productId }">
                    <div class="card-body text-center">
                      <h5 class="card-title">{{ prod.product.name }}</h5>
                      <p class="card-text">{{ $t('price') }}: <strong>{{ prod.product.price }} {{ $t('currency') }}</strong></p>
                      <p class="card-text">{{ $t('duration') }}: <strong>{{ prod.product.description }} {{ $t('days') }}</strong></p>
                      <button class="btn btn-outline-primary w-100" @click="selectProduct(levelId, prod.product)">
                        {{ selectedLevel === levelId && selectedProduct?.id === prod.productId ? $t('selected') : $t('choose') }}
                      </button>
                    </div>
                  </div>
                </div>
                <div v-if="!products || products.length === 0" class="col-12 text-center text-muted">
                  {{ ($t('no_products_for_level') || 'Brak produktów dla tego poziomu.').toLowerCase() }}
                </div>
              </div>
            </div>
            <h4 class="h5 mt-4 mb-3">{{ $t('choose_payment_method') }}</h4>
            <div class="d-flex gap-3 justify-content-center mb-4">
              <button v-for="method in paymentMethods" :key="method.key" class="btn btn-outline-secondary" :class="{active: paymentMethod === method.key}" @click="paymentMethod = method.key">
                {{ method.label }}
              </button>
            </div>
            <div class="text-center">
              <button class="btn btn-success px-5" :disabled="!selectedLevel || !selectedProduct || !paymentMethod" @click="handleSubscribe">
                {{ $t('proceed_to_payment') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { SubscriptionLevelProduct } from '@/types/subscription_level_product';
import type { Product } from '~/types/product';
import type { OrderCreatedResponse } from '~/types/order_created';

definePageMeta({
    middleware: ['auth'],
    layout: 'panel',
})

const subscriptionLevels = ref<Map<number, SubscriptionLevelProduct[]>>(new Map());
const api = useApi();
const selectedLevel = ref<number|null>(null);
const selectedProduct = ref<Product|null>(null);
const paymentMethod = ref<string|null>(null);

const selectedLevelName = computed(() => {
  if (!selectedLevel.value) return '';
  const products = subscriptionLevels.value.get(selectedLevel.value);
  return products && products[0]?.subscriptionLevel?.name ? products[0].subscriptionLevel.name : '';
});

const selectedProductName = computed(() => {
  if (!selectedLevel.value || !selectedProduct.value) return '';
  const products = subscriptionLevels.value.get(selectedLevel.value);
  const prod = products?.find(p => p.productId === selectedProduct.value?.id);
  return prod ? prod.product.name : '';
});
const paymentMethodName = computed(() => {
  const found = paymentMethods.value.find(m => m.key === paymentMethod.value);
  return found ? found.label : '';
});

const paymentMethods = ref<{ key: string, label: string }[]>([])

function selectProduct(levelId: number, productId: Product) {
  selectedLevel.value = levelId;
  selectedProduct.value = productId;
}

async function handleSubscribe() {
  const response = await api.post({
    url: '/order/',
    body: {
      items: [
        {
          product_id: selectedProduct.value?.id,
          quantity: 1
        },
      ],
      payment_method: paymentMethod.value,
      total_amount: selectedProduct.value?.price,
    }
  })
  if (response.statusCode === 201 && response.body?.success) {
    // showSuccess.value = true;
    alert($t('subscription_successful'));
    const data = response.body.data as OrderCreatedResponse;
    if (data.status === 'pending') {
      setTimeout(() => {
        navigateTo(data.redirect_url, { external: true });
      }, 1000);
    }
  } else {
    alert($t('subscription_failed') || 'Subskrypcja nie powiodła się. Spróbuj ponownie.');
  }
}

api.get({url:'/user/subscription_levels'}).then(res => {
  if (res.statusCode === 200 && res.body?.success && res.body.data) {
    const entries: [number, SubscriptionLevelProduct[]][] = Object.entries(res.body.data)
      .filter(([_, v]) => Array.isArray(v))
      .map(([k, v]) => [Number(k),v as SubscriptionLevelProduct[]]);
    subscriptionLevels.value = new Map(entries);
  }
}).catch(err => {
  console.error('Failed to load subscription levels:', err);
});

api.get({url:'/order/payment_methods'}).then(res => {
  if (res.statusCode === 200 && res.body?.success && Array.isArray(res.body.data)) {
    res.body.data.forEach((v : string) => {
      paymentMethods.value.push({
        key: v,
        label: $t('payment_method_' + v)
      });
    })
  }
}).catch(err => {
  console.error('Failed to load payment methods:', err);
});
</script>

<style scoped>
.card.border-primary {
  border-width: 2px;
}
.btn.active {
  background: #0d6efd;
  color: #fff;
}
</style>
