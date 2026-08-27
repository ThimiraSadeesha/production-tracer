DROP PROCEDURE IF EXISTS shift_roster_save;
CREATE PROCEDURE shift_roster_save(
    IN p_date DATETIME,
    IN p_shift_id BIGINT,
    IN p_assignments JSON,
    IN p_created_by VARCHAR(255)
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE total_assignments INT DEFAULT 0;
    DECLARE p_operator_id BIGINT;
    DECLARE p_machine_id BIGINT;

    IF p_assignments IS NULL OR JSON_LENGTH(p_assignments) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Assignments JSON cannot be empty';
    END IF;

    SET total_assignments = JSON_LENGTH(p_assignments);

    SET i = 0;
    WHILE i < total_assignments
        DO
            SET p_operator_id =
                    CAST(JSON_UNQUOTE(JSON_EXTRACT(p_assignments, CONCAT('$[', i, '].operatorId'))) AS UNSIGNED);

            SET p_machine_id = CAST(
                    NULLIF(JSON_UNQUOTE(JSON_EXTRACT(p_assignments, CONCAT('$[', i, '].machineId'))), 'null')
                AS UNSIGNED
                               );
            IF p_operator_id IS NULL THEN
                SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Operator ID cannot be NULL';
            END IF;
            IF EXISTS (SELECT 1
                       FROM tbl_shift_roster
                       WHERE DATE(date) = DATE(p_date)
                         AND shift_id = p_shift_id
                         AND operator_id = p_operator_id
                         AND status = 'Active'
                         AND (machine_id = p_machine_id OR (p_machine_id IS NULL AND machine_id IS NULL))) THEN

                UPDATE tbl_shift_roster
                SET status     = 'Inactive',
                    updated_at = CURRENT_TIMESTAMP(3),
                    updated_by = p_created_by
                WHERE DATE(date) = DATE(p_date)
                  AND shift_id = p_shift_id
                  AND operator_id = p_operator_id
                  AND (machine_id = p_machine_id OR (p_machine_id IS NULL AND machine_id IS NULL));

            END IF;
            INSERT INTO tbl_shift_roster (date, shift_id, machine_id, operator_id, status, created_by)
            VALUES (p_date, p_shift_id, p_machine_id, p_operator_id, 'Active', p_created_by);

            SET i = i + 1;
        END WHILE;

    Select total_assignments As SavedAssignments;

END;