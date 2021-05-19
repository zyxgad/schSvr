
package kpnmsql

import(
	sql   "database/sql"

	mysql "github.com/go-sql-driver/mysql"
	util  "github.com/zyxgad/go-util/util"
)


type SqlTable interface{
	Name()(string)

	SqlInsert(obj Map)(ok bool)

	SqlDelete(cond WhereMap)(ok bool)

	SqlUpdate(obj Map, cond WhereMap)(ok bool)

	SqlSearch(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map)
}

type Table struct{
	SqlTable

	name string
}

func NewSqlTable(name string)(tb *Table){
	return &Table{
		name: name,
	}
}

func (tb *Table)Name()(string){
	return tb.name
}

func (tb *Table)sqlInsert(keys []string, values []interface{})(ok bool){
	var(
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
	)
	tx, err = sqldb.Begin()
	util.AssertError(err, "Sql begin error: ")
	defer func(){ if tx != nil { tx.Rollback() } }()

	order := FormatValue(keys)

	stmt, err = tx.Prepare(util.FormatStr("INSERT INTO %s %s", tb.name, order))
	util.AssertError(err, "Create stmt error: ")
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		switch verr := err.(type) {
		case *mysql.MySQLError:
			switch verr.Number {
			case 1062:
				return false
			}
		}
		util.AssertError(err, "Sql insert error: ")
		return false
	}
	tx.Commit(); tx = nil

	return true
}

func (tb *Table)SqlInsert(obj Map)(ok bool){
	util.Assert(obj != nil && len(obj) > 0, "Obj can't be nil and it's length must greater than 0")

	keys, values := getMapKValues(obj)
	return tb.sqlInsert(keys, values)
}

func (tb *Table)SqlDelete(cond WhereMap)(ok bool){
	if len(cond) == 0{
		return true
	}

	var(
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
	)
	tx, err = sqldb.Begin()
	util.AssertError(err, "Sql begin error: ")
	defer func(){ if tx != nil { tx.Rollback() } }()

	order, values := FormatWhere(cond)

	stmt, err = tx.Prepare(util.FormatStr("DELETE FROM %s %s", tb.name, order))
	util.AssertError(err, "Create stmt error: ")
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		util.AssertError(err, "Sql delete error: ")
		return false
	}
	tx.Commit(); tx = nil

	return true
}

func (tb *Table)sqlUpdate(keys []string, values []interface{}, cond WhereMap)(ok bool){
	var(
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
	)
	tx, err = sqldb.Begin()
	util.AssertError(err, "Sql begin error: ")
	defer func(){ if tx != nil { tx.Rollback() } }()

	where_order, where_values := FormatWhere(cond)

	order := FormatSet(keys)

	stmt, err = tx.Prepare(util.FormatStr("UPDATE %s %s %s", tb.name, order, where_order))
	util.AssertError(err, "Create stmt error: ")
	defer stmt.Close()

	_, err = stmt.Exec(append(values, where_values...)...)
	if err != nil {
		switch verr := err.(type) {
		case *mysql.MySQLError:
			switch verr.Number {
			case 1069:
				return false
			}
		}
		util.AssertError(err, "Sql update error: ")
		return false
	}
	tx.Commit(); tx = nil

	return true
}

func (tb *Table)SqlUpdate(obj Map, cond WhereMap)(ok bool){
	keys, values := getMapKValues(obj)
	return tb.sqlUpdate(keys, values, cond)
}

func (tb *Table)sqlSearch(keys []string, types []string, cond WhereMap, offset uint, limit uint)(lines []Map){
	var(
		err  error
		tx   *sql.Tx
		stmt *sql.Stmt
		rows *sql.Rows
	)
	tx, err = sqldb.Begin()
	util.AssertError(err, "Sql begin error: ")
	defer func(){ if tx != nil { tx.Rollback() } }()

	where_order, values := FormatWhere(cond)
	order := FormatKeyList(keys)

	if limit == 0 {
		stmt, err = tx.Prepare(util.FormatStr("SELECT %s FROM %s %s", order, tb.name, where_order))
	}else{
		stmt, err = tx.Prepare(util.FormatStr("SELECT %s FROM %s %s LIMIT %d, %d", order, tb.name, where_order, offset, limit))
	}
	util.AssertError(err, "Create stmt error: ")
	defer stmt.Close()

	rows, err = stmt.Query(values...)
	if err != nil {
		panic(err)
		return nil
	}
	defer rows.Close()

	lines = make([]Map, 0)
	for rows.Next() {
		rowa := makeScanRowByTypes(types)
		err = rows.Scan(rowa...)
		if err != nil {
			util.AssertError(err, "Sql search error: ")
			return nil
		}
		lines = append(lines, BindKeyValue(keys, rowa))
	}
	tx.Commit(); tx = nil

	return lines
}


func (tb *Table)SqlSearch(obj TypeMap, cond WhereMap, offset uint, limit uint)(lines []Map){
	keys, types := SplitKeyType(obj)
	return tb.sqlSearch(keys, types, cond, offset, limit)
}


