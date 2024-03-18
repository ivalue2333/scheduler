-- noinspection SqlNoDataSourceInspectionForFile

-- noinspection SqlDialectInspectionForFile

CREATE TABLE `scheduler_instance`
(
    `id`            bigint(20) unsigned NOT NULL COMMENT '主键',
    `heart_beat_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '心跳时间',
    `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY             `idx_heart_beat_at` (`heart_beat_at`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='调度实例表';


CREATE TABLE `scheduler_locker`
(
    `id`          bigint(20) unsigned NOT NULL COMMENT '主键',
    `key`         char(64) NOT NULL COMMENT '存储的 key',
    `instance_id` bigint(20) unsigned NOT NULL COMMENT 'instance 主键',
    `created_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `udx_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='调度 locker 存储';

--  drop table scheduler_task;
CREATE TABLE `scheduler_task`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `instance_id` bigint(20) unsigned NOT NULL COMMENT '执行实例',
    `name`        varchar(256) NOT NUll COMMENT '任务名,绑定执行器',
    `namespace`   varchar(64)  NOT NUll COMMENT '实例空间,c0',
    `category`    varchar(64)  NOT NULL DEFAULT '' COMMENT '实例空间分类,c1,推荐拼接',
    `priority`    int          NOT NULL DEFAULT 100 COMMENT '任务优先级, 1 - 100, 1:数值越小,优先级越高',
    `status`      int unsigned NOT NULL DEFAULT 1 COMMENT '状态; 1:等待中,2:执行中,3:暂停,4:停止,5:成功,6:过期',
    `extra`       varchar(1024) COMMENT '业务数据',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_namespace_category_status` (`namespace`, `category`, `status`),
    KEY           `idx_instance_id_status_id` (`instance_id`,`status`, `id`)
) ENGINE=InnoDB  AUTO_INCREMENT = 1 DEFAULT CHARSET=utf8mb4 COMMENT='任务表';