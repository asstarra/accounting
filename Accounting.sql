SET GLOBAL FOREIGN_KEY_CHECKS=1;

USE mysql;
DROP DATABASE a;
USE mysql;
CREATE DATABASE a ; 
USE a;

-- Основные таблицы.

CREATE TABLE entity_type (	
	id SMALLINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE);

INSERT entity_type (id, title) VALUES (1, 'Заказ'), (2, 'Комплект'), (3, 'Модуль'), (4, 'Корпус'), (5, 'Узел');

CREATE TABLE entity (
	id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	id_type SMALLINT NOT NULL,
	specification VARCHAR(255) NOT NULL,
	marking TINYINT NOT NULL DEFAULT 0 CHECK(marking >= 0 AND marking < 4),
	enumerable BOOL NOT NULL DEFAULT true,
	note VARCHAR(1023) NOT NULL,
	UNIQUE (id_type, title),
	FOREIGN KEY (id_type) REFERENCES entity_type(id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE entity_rec (
	id_parent BIGINT NOT NULL,
	id_child BIGINT NOT NULL,
	count INT NOT NULL,
	PRIMARY KEY (id_parent, id_child),
	FOREIGN KEY (id_parent) REFERENCES entity (id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_child) REFERENCES entity (id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE marking (
	id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY);

CREATE TABLE marking_line (
	id_marking BIGINT NOT NULL,
	id_entity BIGINT NOT NULL,
	number TINYINT NOT NULL,
	UNIQUE (id_marking, number),
	PRIMARY KEY (id_marking, id_entity),
	FOREIGN KEY (id_marking) REFERENCES marking(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_entity) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE marked_detail (
	id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	id_marking BIGINT NOT NULL,
	mark VARCHAR(15) NOT NULL,
	id_parent BIGINT NULL,
	UNIQUE(id_marking, mark),
	FOREIGN KEY (id_marking) REFERENCES marking(id) ON DELETE RESTRICT ON UPDATE CASCADE, -- CASCADE?
	FOREIGN KEY (id_parent) REFERENCES marked_detail(id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE person (
	id SMALLINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	name VARCHAR(63) NOT NULL UNIQUE);

CREATE TABLE status_type (	
	id SMALLINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE);
	
CREATE TABLE status (	
	id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	day DATE NOT NULL,
	id_type SMALLINT NOT NULL,
	id_person SMALLINT NULL,
	note VARCHAR(63) NOT NULL,
	FOREIGN KEY (id_person) REFERENCES person(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_type) REFERENCES status_type(id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE status_detail (
	id_detail BIGINT NOT NULL,
	id_status BIGINT NOT NULL,
	PRIMARY KEY (id_detail, id_status),
	FOREIGN KEY (id_detail) REFERENCES marked_detail(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_status) REFERENCES status(id) ON DELETE CASCADE ON UPDATE CASCADE);

CREATE TABLE operation (	
	id SMALLINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE);

CREATE TABLE qualification (
	id_person SMALLINT NOT NULL,
	id_operation SMALLINT NOT NULL,
	level TINYINT NOT NULL,
	PRIMARY KEY (id_person, id_operation),
	FOREIGN KEY (id_person) REFERENCES person(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_operation) REFERENCES operation(id) ON DELETE CASCADE ON UPDATE CASCADE);

CREATE TABLE route_sheet (
	id_entity BIGINT NOT NULL,
	number SMALLINT NOT NULL,
	duration INT NOT NULL,
	person_count TINYINT NOT NULL,
	id_operation SMALLINT NOT NULL,
	PRIMARY KEY (id_entity, number),
	FOREIGN KEY (id_entity) REFERENCES entity(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_operation) REFERENCES operation(id) ON DELETE RESTRICT ON UPDATE CASCADE);

CREATE TABLE detail (
	id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	id_entity BIGINT NOT NULL,
	state TINYINT NOT NULL,
	start_time DATETIME NOT NULL,
	finish_time DATETIME NOT NULL,
	id_parent BIGINT NULL,
	FOREIGN KEY (id_entity) REFERENCES entity(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY (id_parent) REFERENCES detail(id) ON DELETE CASCADE ON UPDATE CASCADE);

CREATE TABLE person_time (
	id_person SMALLINT NOT NULL,
	start_time DATETIME NOT NULL,
	finish_time DATETIME NOT NULL,
	id_detail BIGINT NULL,
	id_entity BIGINT NULL,
	number SMALLINT NULL,
	UNIQUE(id_person, start_time),
	FOREIGN KEY (id_person) REFERENCES person(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_entity, number) REFERENCES route_sheet(id_entity, number) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_detail) REFERENCES detail(id) ON DELETE RESTRICT ON UPDATE CASCADE); --CASCADE


-- Дополнительные таблицы и представления.

-- First and Last element in Marking_Line (id_marking, id_order, id_max_elem, max_elem_number)
CREATE VIEW FLML AS SELECT f.id_marking, f.id_entity AS id_order, l.id_entity AS id_max_elem, l.number AS max_elem_number
	FROM (SELECT id_marking AS id_marking, MAX(number) AS max_number FROM marking_line GROUP BY id_marking) AS n,
		marking_line AS f, marking_line AS l 
	WHERE f.id_marking = l.id_marking AND f.number = 1 AND l.id_marking = n.id_marking AND l.number = n.max_number;
-- two elements side by side IN Marking_Line with condition
CREATE TABLE T_INML (
	id_marking BIGINT NOT NULL,
	id_parent BIGINT NOT NULL,
	parent_number TINYINT NOT NULL,
	id_child BIGINT NOT NULL,
	child_number TINYINT NOT NULL,
	PRIMARY KEY (id_marking, id_parent, id_child),
	FOREIGN KEY (id_parent) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_child) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE);
-- Order and Max element Entity id from T_INML
CREATE TABLE T_OME (
	id_order BIGINT NOT NULL,
	id_max_elem BIGINT NOT NULL,
	FOREIGN KEY (id_order) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_max_elem) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE);
-- First and Last element in Marking_Line with condition
CREATE TABLE T_FLML (
	id_marking BIGINT NOT NULL PRIMARY KEY,
	id_order BIGINT NOT NULL,
	id_max_elem BIGINT NOT NULL,
	max_elem_number TINYINT NOT NULL,
	FOREIGN KEY (id_order) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_max_elem) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE);
CREATE TABLE T_WAY (
	id_marking BIGINT NOT NULL PRIMARY KEY,
	id_order BIGINT NOT NULL,
	id_max_elem BIGINT NOT NULL,
	way_count TINYINT NOT NULL,
	FOREIGN KEY (id_order) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	FOREIGN KEY (id_max_elem) REFERENCES entity(id) ON DELETE RESTRICT ON UPDATE CASCADE);



 DELIMITER $$
CREATE TRIGGER entity_cascade_delete BEFORE DELETE ON entity FOR EACH ROW
BEGIN
	DELETE FROM marking WHERE id IN (SELECT id_marking FROM marking_line WHERE id_entity = OLD.id);
END $$

CREATE TRIGGER entity_rec_cascade_delete BEFORE DELETE ON entity_rec FOR EACH ROW
BEGIN
	INSERT INTO T_INML (id_marking, id_parent, parent_number, id_child, child_number) 
		SELECT p.id_marking, p.id_entity, p.number, c.id_entity, c.number 
			FROM marking_line AS c, marking_line AS p 
			WHERE c.id_marking = p.id_marking AND c.number = p.number + 1
				AND p.id_entity = OLD.id_parent AND c.id_entity = OLD.id_child;
	INSERT INTO T_OME (id_order, id_max_elem) SELECT f.id_order, f.id_max_elem
			FROM FLML AS f, T_INML AS e WHERE f.id_marking = e.id_marking;
	INSERT INTO T_FLML (id_marking, id_order, id_max_elem, max_elem_number) 
		SELECT id_marking, id_order, id_max_elem, max_elem_number FROM FLML 
			WHERE id_order IN (SELECT id_order FROM T_OME) 
				AND id_max_elem IN (SELECT id_max_elem FROM T_OME);
	INSERT INTO T_WAY (id_marking, id_order, id_max_elem, way_count)
		SELECT f.id_marking, f.id_order, f.id_max_elem, w.way_count 
			FROM T_FLML AS f,
				(SELECT id_order, id_max_elem, COUNT(*) AS way_count 
					FROM (SELECT id_order, id_max_elem FROM T_FLML
						UNION ALL
						SELECT id_parent, id_child FROM entity_rec 
							WHERE id_parent IN (SELECT id_order FROM T_OME) 
								AND id_child IN (SELECT id_max_elem FROM T_OME)) AS r
					GROUP BY id_order, id_max_elem) AS w
			WHERE f.id_order = w.id_order AND f.id_max_elem = w.id_max_elem AND f.max_elem_number = 2;
	DELETE FROM marking WHERE id IN 
		(SELECT id_marking FROM T_INML WHERE id_marking NOT IN (SELECT id_marking FROM T_WAY WHERE way_count >= 3)
		UNION
		SELECT id_marking FROM T_WAY WHERE way_count < 3);
	DELETE FROM T_INML;
	DELETE FROM T_OME;
	DELETE FROM T_FLML;
	DELETE FROM T_WAY;
END $$
 DELIMITER ;

