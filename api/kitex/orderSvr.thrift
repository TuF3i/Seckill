namespace go ordersvr

enum OrderStatus {
    UNPAID = 1
    PAID = 2
    CANCELLED = 3
}

struct OrderInfo {
    1: required string orderId
    2: required string userId
    3: required string itemId
    4: required double price
    5: required OrderStatus status
    6: required string createTime
}

service OrderSvr {
    string CreateOrder(1: string userId, 2: string itemId, 3: double price)
    list<OrderInfo> QueryPaidOrders(1: string userId)
    list<OrderInfo> QueryUnpaidOrders(1: string userId)
    list<OrderInfo> QueryCancelledOrders(1: string userId)
}
