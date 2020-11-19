/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50709
 Source Host           : 127.0.0.1:3306
 Source Schema         : go-dog-ctl

 Target Server Type    : MySQL
 Target Server Version : 50709
 File Encoding         : 65001

 Date: 19/11/2020 18:22:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin`  (
  `admin_id` bigint(20) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `phone` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `pwd` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `salt` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `level` int(11) NULL DEFAULT NULL,
  `owner_id` bigint(20) NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  `is_online` tinyint(1) NULL DEFAULT NULL,
  `gate_address` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `role_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`admin_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin
-- ----------------------------
INSERT INTO `admin` VALUES (6770985821484855297, 'admin', '13688460148', '172ab738e6351f63629a7eab2602fcb3', 'Tco00R', 1, 6770985821484855296, 0, 0, '', 1, 1605585783);

-- ----------------------------
-- Table structure for build_service
-- ----------------------------
DROP TABLE IF EXISTS `build_service`;
CREATE TABLE `build_service`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `admin_id` bigint(20) NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT NULL,
  `owner_id` bigint(20) NULL DEFAULT NULL,
  `log` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for docker
-- ----------------------------
DROP TABLE IF EXISTS `docker`;
CREATE TABLE `docker`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `admin_id` bigint(20) NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `account` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `pwd` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `ports` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `owner_id` bigint(20) NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log`  (
  `log_id` bigint(20) NOT NULL,
  `type` int(11) NULL DEFAULT NULL,
  `admin_id` bigint(20) NULL DEFAULT NULL,
  `admin_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `method` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `owner_id` bigint(20) NULL DEFAULT NULL,
  `ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`log_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of log
-- ----------------------------
INSERT INTO `log` VALUES (6770699347048771584, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605664259);
INSERT INTO `log` VALUES (6770699347048771585, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605664267);
INSERT INTO `log` VALUES (6770699347048771586, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605670836);
INSERT INTO `log` VALUES (6770699347048771587, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605670913);
INSERT INTO `log` VALUES (6770699347048771588, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605670964);
INSERT INTO `log` VALUES (6770699347048771589, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605671483);
INSERT INTO `log` VALUES (6770699347048771590, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605671679);
INSERT INTO `log` VALUES (6770699347048771591, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605672094);
INSERT INTO `log` VALUES (6770699347048771592, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605672963);
INSERT INTO `log` VALUES (6770699347048771593, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605673175);
INSERT INTO `log` VALUES (6770699347048771594, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605673483);
INSERT INTO `log` VALUES (6770699347048771595, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605673638);
INSERT INTO `log` VALUES (6770699347048771596, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605673679);
INSERT INTO `log` VALUES (6770699347048771597, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605674396);
INSERT INTO `log` VALUES (6770699347048771598, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605674436);
INSERT INTO `log` VALUES (6770699347048771599, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605674442);
INSERT INTO `log` VALUES (6770839255608438784, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605684499);
INSERT INTO `log` VALUES (6770839586320920576, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605680916);
INSERT INTO `log` VALUES (6770839590666219520, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605680950);
INSERT INTO `log` VALUES (6770839792647122944, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605685160);
INSERT INTO `log` VALUES (6770840415199277056, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605681139);
INSERT INTO `log` VALUES (6770840415199277057, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605681478);
INSERT INTO `log` VALUES (6770840415199277058, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605681512);
INSERT INTO `log` VALUES (6770840419695570944, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605681115);
INSERT INTO `log` VALUES (6770840424007315456, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605682713);
INSERT INTO `log` VALUES (6770840707408048128, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605685370);
INSERT INTO `log` VALUES (6770843640753270784, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605685553);
INSERT INTO `log` VALUES (6770844529979273216, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605681571);
INSERT INTO `log` VALUES (6770844529979273217, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605681580);
INSERT INTO `log` VALUES (6770844529979273218, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605681597);
INSERT INTO `log` VALUES (6770844529979273219, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605681619);
INSERT INTO `log` VALUES (6770844529979273220, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605682314);
INSERT INTO `log` VALUES (6770844529979273221, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605682332);
INSERT INTO `log` VALUES (6770844529979273222, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605682353);
INSERT INTO `log` VALUES (6770844529979273223, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605682358);
INSERT INTO `log` VALUES (6770844529979273224, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605682370);
INSERT INTO `log` VALUES (6770844529979273225, 5, 6770985821484855297, 'admin', 'DelMenu', '删除菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/menu', 1605682395);
INSERT INTO `log` VALUES (6770844529979273226, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605682440);
INSERT INTO `log` VALUES (6770844529979273227, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605682447);
INSERT INTO `log` VALUES (6770845380198248448, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605684395);
INSERT INTO `log` VALUES (6770845380198248449, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605684409);
INSERT INTO `log` VALUES (6770874431390593024, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605688733);
INSERT INTO `log` VALUES (6770874431390593025, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605689041);
INSERT INTO `log` VALUES (6770874431390593026, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605689077);
INSERT INTO `log` VALUES (6770874431390593027, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605689092);
INSERT INTO `log` VALUES (6770874431390593028, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605689100);
INSERT INTO `log` VALUES (6770874431390593029, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605690241);
INSERT INTO `log` VALUES (6770874431390593030, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605690319);
INSERT INTO `log` VALUES (6770875260336058368, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605691020);
INSERT INTO `log` VALUES (6770878563300126720, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605691774);
INSERT INTO `log` VALUES (6770878816602533888, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605693940);
INSERT INTO `log` VALUES (6770878816602533889, 6, 6770985821484855297, 'admin', 'CreateRole', '创建角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/role', 1605694409);
INSERT INTO `log` VALUES (6770878816602533890, 6, 6770985821484855297, 'admin', 'CreateRole', '创建角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/role', 1605694470);
INSERT INTO `log` VALUES (6770879989061496832, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605692155);
INSERT INTO `log` VALUES (6770880199582003200, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605692205);
INSERT INTO `log` VALUES (6770880268335034368, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605692180);
INSERT INTO `log` VALUES (6770985821484855298, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605586032);
INSERT INTO `log` VALUES (6770985821484855299, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605586036);
INSERT INTO `log` VALUES (6771020125287985152, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605591316);
INSERT INTO `log` VALUES (6771020378623946752, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605591901);
INSERT INTO `log` VALUES (6771020378623946753, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592018);
INSERT INTO `log` VALUES (6771020378623946754, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592069);
INSERT INTO `log` VALUES (6771020378623946755, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592155);
INSERT INTO `log` VALUES (6771020378623946756, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592165);
INSERT INTO `log` VALUES (6771020378623946757, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592175);
INSERT INTO `log` VALUES (6771020378623946758, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592286);
INSERT INTO `log` VALUES (6771020378623946759, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592358);
INSERT INTO `log` VALUES (6771020378623946760, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592397);
INSERT INTO `log` VALUES (6771020378623946761, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592404);
INSERT INTO `log` VALUES (6771020378623946762, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592831);
INSERT INTO `log` VALUES (6771020378623946763, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605592948);
INSERT INTO `log` VALUES (6771020378623946764, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605593223);
INSERT INTO `log` VALUES (6771020378623946765, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605593927);
INSERT INTO `log` VALUES (6771020378623946766, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594029);
INSERT INTO `log` VALUES (6771020378623946767, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594039);
INSERT INTO `log` VALUES (6771020378623946768, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594361);
INSERT INTO `log` VALUES (6771020378623946769, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594435);
INSERT INTO `log` VALUES (6771020378623946770, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594468);
INSERT INTO `log` VALUES (6771020378623946771, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594640);
INSERT INTO `log` VALUES (6771020378623946772, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594664);
INSERT INTO `log` VALUES (6771020378623946773, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594698);
INSERT INTO `log` VALUES (6771020378623946774, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594743);
INSERT INTO `log` VALUES (6771020378623946775, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594863);
INSERT INTO `log` VALUES (6771020378623946776, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594887);
INSERT INTO `log` VALUES (6771020378623946777, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594934);
INSERT INTO `log` VALUES (6771020378623946778, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605594958);
INSERT INTO `log` VALUES (6771020378623946779, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595055);
INSERT INTO `log` VALUES (6771020378623946780, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595091);
INSERT INTO `log` VALUES (6771020378623946781, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595103);
INSERT INTO `log` VALUES (6771020378623946782, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595144);
INSERT INTO `log` VALUES (6771020378623946783, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595158);
INSERT INTO `log` VALUES (6771020378623946784, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605595223);
INSERT INTO `log` VALUES (6771020378623946785, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598635);
INSERT INTO `log` VALUES (6771020378623946786, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598663);
INSERT INTO `log` VALUES (6771020378623946787, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598670);
INSERT INTO `log` VALUES (6771020378623946788, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598683);
INSERT INTO `log` VALUES (6771020378623946789, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598724);
INSERT INTO `log` VALUES (6771020378623946790, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598865);
INSERT INTO `log` VALUES (6771020378623946791, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605598906);
INSERT INTO `log` VALUES (6771020378623946792, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599072);
INSERT INTO `log` VALUES (6771020378623946793, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599129);
INSERT INTO `log` VALUES (6771020378623946794, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599149);
INSERT INTO `log` VALUES (6771020378623946795, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599166);
INSERT INTO `log` VALUES (6771020378623946796, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599203);
INSERT INTO `log` VALUES (6771020378623946797, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599267);
INSERT INTO `log` VALUES (6771020378623946798, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599442);
INSERT INTO `log` VALUES (6771020378623946799, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599556);
INSERT INTO `log` VALUES (6771020378623946800, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599881);
INSERT INTO `log` VALUES (6771020378623946801, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605599938);
INSERT INTO `log` VALUES (6771020378623946802, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605600085);
INSERT INTO `log` VALUES (6771020378623946803, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605600171);
INSERT INTO `log` VALUES (6771020378623946804, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605600531);
INSERT INTO `log` VALUES (6771020378623946805, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605600854);
INSERT INTO `log` VALUES (6771020378623946806, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605601239);
INSERT INTO `log` VALUES (6771020378623946807, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605602195);
INSERT INTO `log` VALUES (6771020378623946808, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605602606);
INSERT INTO `log` VALUES (6771020378623946809, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605602701);
INSERT INTO `log` VALUES (6771020378623946810, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605604544);
INSERT INTO `log` VALUES (6771020378623946811, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605604918);
INSERT INTO `log` VALUES (6771020378623946812, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605605181);
INSERT INTO `log` VALUES (6771020378623946813, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605605375);
INSERT INTO `log` VALUES (6771020455950135296, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605591419);
INSERT INTO `log` VALUES (6771020455950135297, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605591693);
INSERT INTO `log` VALUES (6771120777846566912, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605751660);
INSERT INTO `log` VALUES (6771120777846566913, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605753232);
INSERT INTO `log` VALUES (6771120777846566914, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605762118);
INSERT INTO `log` VALUES (6771120777846566915, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605762273);
INSERT INTO `log` VALUES (6771120777846566916, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605763430);
INSERT INTO `log` VALUES (6771122452867035136, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605749963);
INSERT INTO `log` VALUES (6771122452867035137, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605750031);
INSERT INTO `log` VALUES (6771122452867035138, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605750261);
INSERT INTO `log` VALUES (6771122452867035139, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605750584);
INSERT INTO `log` VALUES (6771124926684311552, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605765748);
INSERT INTO `log` VALUES (6771124926684311553, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role', 1605765767);
INSERT INTO `log` VALUES (6771125450687098880, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605765841);
INSERT INTO `log` VALUES (6771125450687098881, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role', 1605765844);
INSERT INTO `log` VALUES (6771125450687098882, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role', 1605765848);
INSERT INTO `log` VALUES (6771125450687098883, 6, 6770985821484855297, 'admin', 'CreateRole', '创建角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/role', 1605766681);
INSERT INTO `log` VALUES (6771125450687098884, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605766685);
INSERT INTO `log` VALUES (6771125450687098885, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605766749);
INSERT INTO `log` VALUES (6771125450687098886, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605766755);
INSERT INTO `log` VALUES (6771125450687098887, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605766762);
INSERT INTO `log` VALUES (6771125450687098888, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605766771);
INSERT INTO `log` VALUES (6771125450687098889, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role/menu', 1605772383);
INSERT INTO `log` VALUES (6771125450687098890, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772388);
INSERT INTO `log` VALUES (6771125450687098891, 6, 6770985821484855297, 'admin', 'CreateRole', '创建角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/role', 1605772422);
INSERT INTO `log` VALUES (6771125450687098892, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772432);
INSERT INTO `log` VALUES (6771125450687098893, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772438);
INSERT INTO `log` VALUES (6771125450687098894, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772447);
INSERT INTO `log` VALUES (6771125450687098895, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772454);
INSERT INTO `log` VALUES (6771125450687098896, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772461);
INSERT INTO `log` VALUES (6771125450687098897, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772472);
INSERT INTO `log` VALUES (6771125450687098898, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605772479);
INSERT INTO `log` VALUES (6771125450687098899, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role/menu', 1605772489);
INSERT INTO `log` VALUES (6771125450687098900, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role/menu', 1605772491);
INSERT INTO `log` VALUES (6771125450687098901, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role/menu', 1605772492);
INSERT INTO `log` VALUES (6771156189868699648, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605774751);
INSERT INTO `log` VALUES (6771160304531255296, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605775756);
INSERT INTO `log` VALUES (6771160304531255297, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605775792);
INSERT INTO `log` VALUES (6771160304531255298, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605775831);
INSERT INTO `log` VALUES (6771160304531255299, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605777967);
INSERT INTO `log` VALUES (6771160304531255300, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778125);
INSERT INTO `log` VALUES (6771160312936640512, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605763708);
INSERT INTO `log` VALUES (6771160312936640513, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605764237);
INSERT INTO `log` VALUES (6771160312936640514, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605764264);
INSERT INTO `log` VALUES (6771160312936640515, 7, 6770985821484855297, 'admin', 'DelRole', '删除角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/role/menu', 1605765241);
INSERT INTO `log` VALUES (6771160312936640516, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605765248);
INSERT INTO `log` VALUES (6771160312936640517, 6, 6770985821484855297, 'admin', 'CreateRole', '创建角色', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/role', 1605765306);
INSERT INTO `log` VALUES (6771160312936640518, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605765311);
INSERT INTO `log` VALUES (6771162035386298368, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605778266);
INSERT INTO `log` VALUES (6771162035386298369, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778300);
INSERT INTO `log` VALUES (6771162035386298370, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778307);
INSERT INTO `log` VALUES (6771162035386298371, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778355);
INSERT INTO `log` VALUES (6771162035386298372, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778362);
INSERT INTO `log` VALUES (6771162035386298373, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778428);
INSERT INTO `log` VALUES (6771162035386298374, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778432);
INSERT INTO `log` VALUES (6771162035386298375, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778440);
INSERT INTO `log` VALUES (6771162035386298376, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778448);
INSERT INTO `log` VALUES (6771162035386298377, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778460);
INSERT INTO `log` VALUES (6771162035386298378, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778479);
INSERT INTO `log` VALUES (6771162035386298379, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778489);
INSERT INTO `log` VALUES (6771162035386298380, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605778493);
INSERT INTO `log` VALUES (6771162035386298381, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605779172);
INSERT INTO `log` VALUES (6771162035386298382, 6, 6770985821484855297, 'admin', 'BindRoleMenu', '绑定角色菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/menu', 1605779178);
INSERT INTO `log` VALUES (6771162035386298383, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779290);
INSERT INTO `log` VALUES (6771162035386298384, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779296);
INSERT INTO `log` VALUES (6771162035386298385, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779304);
INSERT INTO `log` VALUES (6771162035386298386, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779310);
INSERT INTO `log` VALUES (6771162035386298387, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779333);
INSERT INTO `log` VALUES (6771162035386298388, 11, 6770985821484855297, 'admin', 'BindRoleAPI', '绑定角色API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/bind/role/api', 1605779343);
INSERT INTO `log` VALUES (6771162035386298389, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605779438);
INSERT INTO `log` VALUES (6771162035386298390, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780194);
INSERT INTO `log` VALUES (6771162035386298391, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780196);
INSERT INTO `log` VALUES (6771162035386298392, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780199);
INSERT INTO `log` VALUES (6771162035386298393, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780209);
INSERT INTO `log` VALUES (6771162035386298394, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780212);
INSERT INTO `log` VALUES (6771162035386298395, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780215);
INSERT INTO `log` VALUES (6771162035386298396, 10, 6770985821484855297, 'admin', 'DelAPI', '删除API', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/del/api', 1605780223);
INSERT INTO `log` VALUES (6771162035386298397, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605780591);
INSERT INTO `log` VALUES (6771162035386298398, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605780908);
INSERT INTO `log` VALUES (6771162035386298399, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605781014);
INSERT INTO `log` VALUES (6771162035386298400, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605781052);
INSERT INTO `log` VALUES (6771162035386298401, 0, 6770985821484855297, 'admin', 'AdminLogin', '管理员登录', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/admin/login', 1605781087);
INSERT INTO `log` VALUES (6771162035386298402, 4, 6770985821484855297, 'admin', 'CreateMenu', '创建菜单', 6770985821484855296, '127.0.0.1', '/api/go-dog/controller/v1/create/menu', 1605781107);

-- ----------------------------
-- Table structure for owner
-- ----------------------------
DROP TABLE IF EXISTS `owner`;
CREATE TABLE `owner`  (
  `owner_id` bigint(20) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `phone` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `level` int(11) NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  `is_admin_owner` tinyint(1) NULL DEFAULT NULL,
  `role_id` int(10) UNSIGNED NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`owner_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of owner
-- ----------------------------
INSERT INTO `owner` VALUES (6770985821484855296, '系统业主', '13688460148', 1, 0, 1, 1, 1605585783);

SET FOREIGN_KEY_CHECKS = 1;
