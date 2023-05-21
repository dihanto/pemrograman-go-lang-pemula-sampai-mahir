package helper

import "database/sql"

func RollbackOrCommit(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		PanifIfError(errorRollback)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanifIfError(errorCommit)
		return
	}
}
