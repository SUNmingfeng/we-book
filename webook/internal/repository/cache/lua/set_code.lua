local key = KEYS[1]
local cntKey = key..":cnt"
local val = ARGV[1]

local ttl  = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- -1 key存在，但没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- -2 key不存在，可以新建
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", key, 600)
    return 0
else
    return -1
end
