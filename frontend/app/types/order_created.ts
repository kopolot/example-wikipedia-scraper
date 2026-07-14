export type OrderCreatedResponse = {
    id: number
    payment_id: string
    redirect_url: string
    status: OrderStatus
}

export type OrderStatus = 'pending' | 'paid' | 'completed' | 'refunded' | 'canceled' | 'failed'