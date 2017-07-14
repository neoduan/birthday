package mysql

import "testing"

func TestInsertWithTx(t *testing.T) {
	var (
		cli *Cli   = getCli()
		sql string = `insert into t_test (title) values ("hello tittle with no close.");`
	)

	creatTable()
	cli.TxBegin()
	cli.Write(sql)
	cli.TxCommit()
	cli.Close()
	dropTable()
}
