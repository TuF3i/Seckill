local flashKey = KEYS[1]
local stockKey = KEYS[2]
local limitKey = KEYS[3]
local limitTtl = tonumber(ARGV[1])

local flashStatus = redis.call('GET', flashKey)
if not flashStatus or tonumber(flashStatus) == 0 then
    return -1
end

local exists = redis.call('EXISTS', limitKey)
if exists == 1 then
    return -2
end

local stock = redis.call('GET', stockKey)
if not stock or tonumber(stock) <= 0 then
    return -3
end

local remaining = redis.call('DECR', stockKey)
if remaining < 0 then
    redis.call('INCR', stockKey)
    return -3
end

redis.call('SET', limitKey, 1, 'EX', limitTtl)

return remaining