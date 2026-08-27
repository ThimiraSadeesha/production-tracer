DROP PROCEDURE IF EXISTS unit_getAll;
CREATE PROCEDURE unit_getAll()
BEGIN
    SELECT tu.id          AS id,
           tu.unit_code   AS unitCode,
           tu.unit_name   AS unitName,
           tu.unit_symbol AS unitSymbol,
           tu.unit_status AS unitStatus

    FROM tbl_unit tu;

END
