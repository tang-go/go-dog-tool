/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50709
 Source Host           : 127.0.0.1:3306
 Source Schema         : go-dog-auth

 Target Server Type    : MySQL
 Target Server Version : 50709
 File Encoding         : 65001

 Date: 19/11/2020 18:22:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_api`;
CREATE TABLE `sys_api`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `organize` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `describe` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `api` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_api
-- ----------------------------
INSERT INTO `sys_api` VALUES (2, 'go-dog', '获取管理员信息', 'api/go-dog/controller/v1/get/admin/info', 1605585782);
INSERT INTO `sys_api` VALUES (3, 'go-dog', '获取角色列表', 'api/go-dog/controller/v1/get/role/list', 1605585783);
INSERT INTO `sys_api` VALUES (4, 'go-dog', '获取编译发布记录', 'api/go-dog/controller/v1/get/build/service/list', 1605585783);
INSERT INTO `sys_api` VALUES (5, 'go-dog', '获取docker运行服务', 'api/go-dog/controller/v1/get/docker/list', 1605585783);
INSERT INTO `sys_api` VALUES (6, 'go-dog', '获取服务列表', 'api/go-dog/controller/v1/get/service/list', 1605585783);
INSERT INTO `sys_api` VALUES (11, 'go-dog', '管理员登录', 'api/go-dog/controller/v1/admin/login', 1605585783);
INSERT INTO `sys_api` VALUES (12, 'go-dog', '编译发布服务', 'api/go-dog/controller/v1/build/service', 1605585783);
INSERT INTO `sys_api` VALUES (13, 'go-dog', 'docker方式启动服务', 'api/go-dog/controller/v1/strat/docker', 1605585783);
INSERT INTO `sys_api` VALUES (14, 'go-dog', '关闭docker服务', 'api/go-dog/controller/v1/clsoe/docker', 1605585783);
INSERT INTO `sys_api` VALUES (15, 'go-dog', '删除docker服务', 'api/go-dog/controller/v1/del/docker', 1605585783);
INSERT INTO `sys_api` VALUES (16, 'go-dog', '重启docker服务', 'api/go-dog/controller/v1/restart/docker', 1605585783);
INSERT INTO `sys_api` VALUES (19, 'go-dog', '获取菜单', 'api/go-dog/controller/v1/get/menu', 1605664244);
INSERT INTO `sys_api` VALUES (20, 'go-dog', '创建菜单', 'api/go-dog/controller/v1/create/menu', 1605664244);
INSERT INTO `sys_api` VALUES (21, 'go-dog', '删除菜单', 'api/go-dog/controller/v1/del/menu', 1605678061);
INSERT INTO `sys_api` VALUES (23, 'go-dog', '获取API列表', 'api/go-dog/controller/v1/get/api/list', 1605688728);
INSERT INTO `sys_api` VALUES (24, 'go-dog', '获取角色菜单列表', 'api/go-dog/controller/v1/get/role/menu/list', 1605688728);
INSERT INTO `sys_api` VALUES (25, 'go-dog', '获取校色api列表', 'api/go-dog/controller/v1/get/role/api/list', 1605688728);
INSERT INTO `sys_api` VALUES (26, 'go-dog', '删除API', 'api/go-dog/controller/v1/del/api', 1605688728);
INSERT INTO `sys_api` VALUES (27, 'go-dog', '创建角色', 'api/go-dog/controller/v1/create/role', 1605688728);
INSERT INTO `sys_api` VALUES (28, 'go-dog', '删除角色', 'api/go-dog/controller/v1/del/role', 1605688729);
INSERT INTO `sys_api` VALUES (29, 'go-dog', '绑定角色菜单', 'api/go-dog/controller/v1/bind/role/menu', 1605688729);
INSERT INTO `sys_api` VALUES (30, 'go-dog', '删除角色菜单', 'api/go-dog/controller/v1/del/role/menu', 1605688729);
INSERT INTO `sys_api` VALUES (31, 'go-dog', '绑定角色API', 'api/go-dog/controller/v1/bind/role/api', 1605688729);
INSERT INTO `sys_api` VALUES (32, 'go-dog', '删除角色API', 'api/go-dog/controller/v1/del/role/api', 1605688729);

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `organize` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `describe` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `parent_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `sort` int(10) UNSIGNED NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, 'go-dog', '权限管理', '/power', 0, 0, 1605585783);
INSERT INTO `sys_menu` VALUES (2, 'go-dog', '菜单管理', '/power/menu', 1, 0, 1605585783);
INSERT INTO `sys_menu` VALUES (4, 'go-dog', '在线服务', '/service', 0, 999, 1605670836);
INSERT INTO `sys_menu` VALUES (5, 'go-dog', 'docker管理', '/docker', 0, 998, 1605671679);
INSERT INTO `sys_menu` VALUES (6, 'go-dog', '发布docker镜像', '/docker/build', 5, 100, 1605672963);
INSERT INTO `sys_menu` VALUES (8, 'go-dog', '通过docker启动服务', '/docker/start', 5, 99, 1605681597);
INSERT INTO `sys_menu` VALUES (11, 'go-dog', '首页', '/index', 0, 1000, 1605682440);
INSERT INTO `sys_menu` VALUES (12, 'go-dog', '角色管理', '/power/role', 1, 100, 1605689077);
INSERT INTO `sys_menu` VALUES (13, 'go-dog', 'API管理', '/power/api', 1, 1, 1605689092);
INSERT INTO `sys_menu` VALUES (14, 'go-dog', '系统管理', '/sys', 0, 1, 1605781052);
INSERT INTO `sys_menu` VALUES (15, 'go-dog', '管理员管理', '/sys/admin', 14, 100, 1605781107);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `organize` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `describe` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `is_super` tinyint(1) NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES (1, 'go-dog', '超级业主', '超级业主', 1, 1605585783);
INSERT INTO `sys_role` VALUES (4, 'go-dog', '前端', '前端API查看', 0, 1605765306);
INSERT INTO `sys_role` VALUES (6, 'go-dog', '后端', '后端开发', 0, 1605766681);
INSERT INTO `sys_role` VALUES (7, 'go-dog', '普通管理员', '普通管理员', 0, 1605772422);

-- ----------------------------
-- Table structure for sys_role_api
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_api`;
CREATE TABLE `sys_role_api`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `api_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_api
-- ----------------------------
INSERT INTO `sys_role_api` VALUES (1, 4, 11, 1605777967);
INSERT INTO `sys_role_api` VALUES (2, 4, 2, 1605778125);
INSERT INTO `sys_role_api` VALUES (3, 6, 11, 1605778300);
INSERT INTO `sys_role_api` VALUES (4, 6, 2, 1605778307);
INSERT INTO `sys_role_api` VALUES (5, 7, 11, 1605778355);
INSERT INTO `sys_role_api` VALUES (6, 7, 2, 1605778362);
INSERT INTO `sys_role_api` VALUES (7, 6, 6, 1605778428);
INSERT INTO `sys_role_api` VALUES (8, 6, 4, 1605778432);
INSERT INTO `sys_role_api` VALUES (9, 6, 5, 1605778440);
INSERT INTO `sys_role_api` VALUES (10, 6, 12, 1605778448);
INSERT INTO `sys_role_api` VALUES (11, 6, 13, 1605778460);
INSERT INTO `sys_role_api` VALUES (12, 6, 14, 1605778479);
INSERT INTO `sys_role_api` VALUES (13, 6, 15, 1605778489);
INSERT INTO `sys_role_api` VALUES (14, 6, 16, 1605778493);
INSERT INTO `sys_role_api` VALUES (15, 7, 3, 1605779290);
INSERT INTO `sys_role_api` VALUES (16, 7, 4, 1605779296);
INSERT INTO `sys_role_api` VALUES (17, 7, 5, 1605779304);
INSERT INTO `sys_role_api` VALUES (18, 7, 6, 1605779310);
INSERT INTO `sys_role_api` VALUES (19, 7, 24, 1605779333);
INSERT INTO `sys_role_api` VALUES (20, 7, 25, 1605779343);

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `menu_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `role_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `add` tinyint(1) NULL DEFAULT NULL,
  `del` tinyint(1) NULL DEFAULT NULL,
  `update` tinyint(1) NULL DEFAULT NULL,
  `select` tinyint(1) NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
INSERT INTO `sys_role_menu` VALUES (5, 11, 6, 1, 1, 1, 1, 1605766685);
INSERT INTO `sys_role_menu` VALUES (6, 4, 6, 1, 1, 1, 1, 1605766749);
INSERT INTO `sys_role_menu` VALUES (7, 5, 6, 1, 1, 1, 1, 1605766755);
INSERT INTO `sys_role_menu` VALUES (8, 6, 6, 1, 1, 1, 1, 1605766762);
INSERT INTO `sys_role_menu` VALUES (9, 8, 6, 1, 1, 1, 1, 1605766771);
INSERT INTO `sys_role_menu` VALUES (10, 11, 4, 0, 0, 0, 1, 1605772388);
INSERT INTO `sys_role_menu` VALUES (11, 11, 7, 0, 0, 0, 1, 1605772432);
INSERT INTO `sys_role_menu` VALUES (12, 4, 7, 0, 0, 0, 1, 1605772438);
INSERT INTO `sys_role_menu` VALUES (16, 1, 7, 0, 0, 0, 1, 1605772471);
INSERT INTO `sys_role_menu` VALUES (17, 12, 7, 0, 0, 0, 1, 1605772479);
INSERT INTO `sys_role_menu` VALUES (18, 5, 7, 0, 0, 0, 1, 1605779172);
INSERT INTO `sys_role_menu` VALUES (19, 6, 7, 0, 0, 0, 1, 1605779178);

SET FOREIGN_KEY_CHECKS = 1;
