/**------------------------------------------------------------**
 * @filename cstring/
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-08-02 16:17
 * @desc     cstring -
 **------------------------------------------------------------**/
package cstring

func BytesToString(b []byte) string {
	return string(b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}
