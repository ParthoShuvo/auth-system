USE `AuthDB`;

-- TEMPORARY SP to insert permissions
DROP PROCEDURE IF EXISTS `temp_permission_sp` ;

DELIMITER ;;
CREATE PROCEDURE `temp_permission_sp`()
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLSTATE '45000' Select 'Duplicate permission';
	CALL sp_insert_permission('GetPost', 'Fetch a post');
    CALL sp_insert_permission('AddPost', 'Insert a post');
    CALL sp_insert_permission('UpdatePost', 'Edit a post');
    CALL sp_insert_permission('DeletePost', 'Delete a post');
END ;;
DELIMITER ;

CALL `temp_permission_sp`;

DROP PROCEDURE IF EXISTS `temp_permission_sp`;

-- TEMPORARY SP to insert roles and its permissions
DROP PROCEDURE IF EXISTS `temp_role_sp` ;


DELIMITER ;;
CREATE PROCEDURE `temp_role_sp`(IN name varchar(64), IN description varchar(512), IN permission varchar(32))
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLSTATE '45000' Select 'Duplicate role permission';
	CALL sp_insert_role(name, description);
    SET @permissionid = (SELECT P.id from Permission AS P where P.name=permission);
    SET @roleid = (SELECT R.id from Role AS R where R.name=name);
    CALL sp_insert_role_permission(@roleid, @permissionid);
END ;;
DELIMITER ;

# Role Admin and its permissions
CALL `temp_role_sp`('Admin', 'Administrative user', 'GetPost');
CALL `temp_role_sp`('Admin', 'Administrative user', 'AddPost');
CALL `temp_role_sp`('Admin', 'Administrative user', 'UpdatePost');
CALL `temp_role_sp`('Admin', 'Administrative user', 'DeletePost');

# Role Author and its permissions
CALL `temp_role_sp`('Author', 'Only read, create and update access', 'GetPost');
CALL `temp_role_sp`('Author', 'Only read, create and update access', 'AddPost');
CALL `temp_role_sp`('Author', 'Only read, create and update access', 'UpdatePost');

# Role Reader and its permissions
CALL `temp_role_sp`('Reader', 'Only read access', 'GetPost');

DROP PROCEDURE IF EXISTS `temp_role_sp` ;

-- TEMPORARY SP to insert users and its roles
DROP PROCEDURE IF EXISTS `temp_user_sp` ;

DELIMITER ;;
CREATE PROCEDURE `temp_user_sp`(IN firstname varchar(64), IN lastname varchar(64), IN login varchar(64),
    IN password varchar(64), IN role varchar(64))
BEGIN
    DECLARE CONTINUE HANDLER FOR SQLSTATE '45000' Select 'Duplicate user role';
	CALL sp_insert_user(firstname, lastname, login, MD5(password), UUID());
    CALL sp_user_verification_assignment(login, 1);
    SET @userid = (SELECT U.id from User AS U where U.login=login);
    SET @roleid = (SELECT R.id from Role AS R where R.name=role);
    CALL sp_insert_user_role(@userid, @roleid);
END ;;
DELIMITER ;

# Users and roles
CALL `temp_user_sp`('Admin', 'User', 'admin.user@testmail.com', '_LaRa08CRoft', 'Admin');
CALL `temp_user_sp`('Author', 'User1', 'author.user1@testmail.com', '_GOllum#!', 'Author');
CALL `temp_user_sp`('Reader', 'User1', 'reader.user1@testmail.com', 'bUfo_MelanOst!ktus', 'reader');

DROP PROCEDURE IF EXISTS `temp_user_sp` ;
