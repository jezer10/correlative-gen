package database

func InitLogs() {

	DB.Exec(`
		CREATE TABLE IF NOT EXISTS LOGS (
			ID SERIAL PRIMARY KEY,
			TYPE VARCHAR(20),
			MESSAGE TEXT,
			DATE TIMESTAMP
		);
	`)

}

func SaveErrorLog(message string) {
	DB.Exec(`
		INSERT INTO LOGS (TYPE, MESSAGE, DATE)
		VALUES ('ERROR', $1, NOW());
	`, message)

}
func SaveInfoLog(message string) {
	DB.Exec(`
		INSERT INTO LOGS (TYPE, MESSAGE, DATE)
		VALUES ('INFO', $1, NOW());
	`, message)

}
