-- 创建 apisix route 表
CREATE TABLE apisix_route (
    id SERIAL PRIMARY KEY,
    route_id VARCHAR NOT NULL,       -- 通过雪花算法生成的 ID
    content VARCHAR NOT NULL,        -- 保存原本 api 请求和返回的结构类型
    content_yaml VARCHAR NOT NULL,   -- apisix config 配置的格式
    type smallint NOT NULL,          -- 0 下线，1 启用
    status smallint NOT NULL,        -- 0 为前后端使用，1 为管理平台注册的 OCR 服务
    create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    update_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 更新时间
    expire_at timestamptz -- 过期时间
);

-- 创建 apisix upstream 表
CREATE TABLE apisix_upstream (
    id SERIAL PRIMARY KEY,
    stream_id VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    content_yaml VARCHAR NOT NULL,   -- apisix config 配置的格式
    type smallint NOT NULL, -- 0 为前后端使用，1 为管理平台注册的 OCR 服务
    create_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    update_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 记录请求日志的表
CREATE TABLE access_log (
    id SERIAL PRIMARY KEY,
    route_id VARCHAR,
    request_id VARCHAR,
    client_addr VARCHAR,
    iso_time VARCHAR,
    timestamp BIGINT,
    datestamp   BIGINT,
    weeklystamp  BIGINT,
    monthstamp  BIGINT,
    yearstamp  BIGINT,
    time_cost FLOAT,
    request_length INT,
    connection VARCHAR,
    connection_requests VARCHAR,
    uri VARCHAR,
    ori_request VARCHAR,
    query_string VARCHAR,
    status INT,
    bytes_sent INT,
    referer VARCHAR,
    user_agent VARCHAR,
    forwarded_for VARCHAR,
    host VARCHAR,
    node VARCHAR,
    upstream VARCHAR
);
