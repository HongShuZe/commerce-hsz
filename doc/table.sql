/*
商品表
 */
CREATE TABLE `product` (
  `id` int(11) not null AUTO_INCREMENT,
  `productName` varchar(255) DEFAULT NULL COMMENT '商品名',
  `productNum` int(11) DEFAULT NULL COMMENT '商品数量',
  `productImage` varchar(255) DEFAULT NULL COMMENT '商品图片地址',
  `productUrl` varchar(255) DEFAULT NULL COMMENT '商品链接',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

create table `order` (
  `id` int(11) not null auto_increment ,
  `userID` int(11) default null ,
  `productID` int(11) default null ,
  `orderStatus` varchar(255) default null ,
  primary key (`id`)
) ENGINE=InnoDB default charset=utf8;