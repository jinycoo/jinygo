/**------------------------------------------------------------**
 * @filename sql/mysql.go
 * @author   jinycoo - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2019/10/15 09:41
 * @desc     sql - mysql adapter db
 **------------------------------------------------------------**/
package sql

func NewMySQL(c *Config) (db *DB) {
	return newDB(_mysql, c)
}
