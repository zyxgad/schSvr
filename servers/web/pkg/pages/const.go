
package kpnmwebpage

import (
	regexp  "regexp"
)

var (
	reg_name  *regexp.Regexp = regexp.MustCompile(`^[A-Za-z_-][0-9A-Za-z_-]{3,31}$`)
	reg_pwd   *regexp.Regexp = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_-]{7,127}$`)
)


func init(){
}
