
package kpnmwebpage

import (

	jwt  "github.com/zyxgad/go-util/jwt"
)

var (
	jwtEncoder jwt.Encoder = jwt.NewAutoEncoder(60 * 60 * 24 * 15, 2048)
)
