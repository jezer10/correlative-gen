/*

Enter custom T-SQL here that would run after SQL Server has started up. 

*/

CREATE DATABASE bgc;
GO

USE bgc;
GO

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'USUARIOS')
BEGIN
    CREATE TABLE USUARIOS ( 
        ID INT IDENTITY(1,1) PRIMARY KEY,
        DNI VARCHAR(20) NULL,
        NOMBRE VARCHAR(80) NULL,
        APELLIDO VARCHAR(100) NULL,
        FECHANACIMIENTO DATE NULL,
        FECHAREGISTRO DATETIME NULL,
        CLIENTE CHAR(10) NULL,
        FLAG BIT NULL,
        NACIONALIDAD VARCHAR(40) NULL,
        STATUS CHAR(25) NULL,
        STATUSDESCRIPTION VARCHAR(40) NULL,
        STATUSNOTE VARCHAR(800) NULL,
        CORRELATIVO VARCHAR(50) NULL,
        SCORE VARCHAR(50) NULL,
        SCORE_DESCRIPTION NVARCHAR(MAX) NULL,
        SCORE_NOTE NVARCHAR(MAX) NULL,
        RESIDENCIA VARCHAR(20) NULL,
        FECHARESPUESTA DATETIME NULL
    );
END;
GO