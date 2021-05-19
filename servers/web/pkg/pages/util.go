
package kpnmwebpage


type JsonError struct{
	Error string `json:"error"`
	ErrorMsg string `json:"errorMessage"`
	Cause string `json:"cause,omitempty"`
}

func CreateJsonError(err string, errMsg string, cause string)(jerr *JsonError){
	jerr = new(JsonError)
	jerr.Error = err
	jerr.ErrorMsg = errMsg
	jerr.Cause = cause
	return jerr
}

