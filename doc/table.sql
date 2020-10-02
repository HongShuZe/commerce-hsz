/*
商品表
 */
CREATE TABLE `tbl_product` (
  `id` int(11) not null AUTO_INCREMENT,
  `seller_id` int(11) default null COMMENT '卖家id',
  `product_name` varchar(255) DEFAULT NULL COMMENT '商品名',
  `product_num` int(11) DEFAULT NULL COMMENT '商品数量',
  `product_image` varchar(255) DEFAULT NULL COMMENT '商品图片地址',
  `product_url` varchar(255) DEFAULT NULL COMMENT '商品链接',
  `product_price` decimal(11,2) DEFAULT NULL COMMENT '商品价格',
  `product_info` varchar(255) DEFAULT NULL COMMENT '商品信息',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间戳',
  `status` int(11) DEFAULT NULL COMMENT '商品状态(0可用/1已删除等状态)',
  `ext2` text COMMENT '备用字段1',
  PRIMARY KEY (`id`),
  KEY `idx_product_name`(`product_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


/*
订单表
 */
create table `tbl_order` (
  `id` int(11) not null auto_increment,
  `user_id` int(11) default null COMMENT '用户id',
  `seller_id` int(11) default null COMMENT '卖家id',
  `product_id` int(11) default null COMMENT '商品id',
  `order_num` int(11) default null COMMENT '购买商品数量',
  `total_price` decimal(11,2) default null COMMENT '总金额',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建订单日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间戳',
  `status` int(11) default null COMMENT '发货状态(0未发货/1已发货等状态)',
  `ext1` text COMMENT '备用字段1',
  primary key (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB default charset=utf8;

/*
用户表
 */
CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `nick_name` varchar(255) DEFAULT NULL COMMENT '昵称',
  `user_name` varchar(255) DEFAULT NULL COMMENT '用户名',
  `pass_word` varchar(255) DEFAULT NULL COMMENT '用户密码',
  `user_balance` int(11) DEFAULT NULL COMMENT '用户余额',
  `user_type` int(11) DEFAULT NULL COMMENT '用户类型(0 买家/1 卖家等状态)',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建用户日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间戳',
  `status` int(11) DEFAULT NULL COMMENT '用户帐号状态(0可用/1已注销等状态)',
  `ext1` text COMMENT '备用字段1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_nickname` (`nick_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;






