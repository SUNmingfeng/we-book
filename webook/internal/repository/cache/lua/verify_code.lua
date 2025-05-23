local key = KEYS[1]
local cntKey = key..":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    --验证次数用尽
    return -1
end

if code == expectedCode then
    redis.call("set", cntKey, 0)
    --验证码正确
    return 0
else
    redis.call("decr", cntKey)
    --验证码错误
    return -2
end




