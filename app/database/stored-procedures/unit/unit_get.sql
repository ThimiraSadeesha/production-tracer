DROP PROCEDURE IF EXISTS unit_get;
CREATE PROCEDURE unit_get(
    IN unit_id_val INT
)
BEGIN

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            GET DIAGNOSTICS CONDITION 1 @sqlstate = RETURNED_SQLSTATE, @errno = MYSQL_ERRNO, @message_text = MESSAGE_TEXT;
            ROLLBACK;
            SELECT CONCAT('Error: [', @sqlstate, '] ', @message_text) AS error_message;
        END;

    START TRANSACTION;

    SELECT tu.id          AS unitId,
           tu.unit_code   AS unitCode,
           tu.unit_name   AS unitName,
           tu.unit_symbol AS unitSymbol,
           tu.unit_status AS unitStatus

    FROM tbl_unit tu
    WHERE tu.id = unit_id_val;
    COMMIT;
END