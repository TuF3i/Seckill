namespace go itemsvr

struct ItemInfo {
    1: required string id
    2: required string name
    3: required i64 stock
    4: required double price
    5: required string description
}

service ItemSvr {
    string AddItem(1: string name, 2: i64 stock, 3: double price, 4: string description)
    void DeleteItem(1: string id)
    void StartFlashSale(1: string itemId)
    void StopFlashSale(1: string itemId)
    list<ItemInfo> ListItems(1: string uid, 2: string role)
}
