namespace go paymentsvr

service PaymentSvr {
    bool ProcessPayment(1: string orderId, 2: string userId)
}
