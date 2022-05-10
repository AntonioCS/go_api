DROP PROCEDURE IF EXISTS AddVenueTables;

DELIMITER //

CREATE PROCEDURE AddVenueTables()
BEGIN
	DECLARE bDone INT DEFAULT 0;
	DECLARE totalTablesPerVenue TINYINT DEFAULT 10;
	DECLARE numTableSeatsMin TINYINT DEFAULT 5;
	DECLARE numTableSeatsMax TINYINT DEFAULT 15;
	DECLARE vId INT;
	DECLARE counter INT;

	DECLARE curs CURSOR FOR SELECT id FROM venue;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET bDone = 1;
		
	OPEN curs;	
	
	REPEAT
		SET counter = 0;
		FETCH curs INTO vId;
		WHILE counter < totalTablesPerVenue DO		
			INSERT INTO `go_api`.`table` (`venue_id`, `size`) VALUES (vId, FLOOR( RAND() * (numTableSeatsMax-numTableSeatsMin) + numTableSeatsMin));
			SET counter = counter + 1;
		END WHILE;
	UNTIL bDone END REPEAT;
	
	CLOSE curs;	
END //

DELIMITER ;

CALL AddVenueTables();
