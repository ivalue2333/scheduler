-- noinspection SqlNoDataSourceInspectionForFile

-- noinspection SqlDialectInspectionForFile

CREATE TABLE `scheduler_instance`
(
    `id`            bigint(20) unsigned NOT NULL COMMENT '主键',
    `heart_beat_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '心跳时间',
    `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='调度实例表';


CREATE TABLE `scheduler_task`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`       varchar(256) NOT NUll COMMENT '任务名,绑定执行器',
    `namespace`  varchar(256) NOT NUll COMMENT '实例空间',
    `category`   varchar(256) NOT NULL DEFAULT '' COMMENT '实例空间分类',
    `priority`   int          NOT NULL DEFAULT 100 COMMENT '任务优先级, 1 - 100, 1:数值越小,优先级越高',
    `status`     int unsigned NOT NULL DEFAULT 1 COMMENT '状态; 1:等待中,2:执行中,3:暂停,4:停止,5:成功,6:过期',
    `extra`      varchar(8192) COMMENT '业务数据',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY          `idx_namespace_category_status` (`namespace`, `category`, `status`)
) ENGINE=InnoDB  AUTO_INCREMENT = 1 DEFAULT CHARSET=utf8mb4 COMMENT='任务表';