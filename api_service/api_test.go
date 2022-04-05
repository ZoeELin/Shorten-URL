package api_service

import (
	"testing"
)


// Test expired date
func TestExpireData(t *testing.T){
	if ExpireData("2010-04-10T10:04Z") == true {
        t.Log("api_service.ExpireData PASS")
    } else {
        t.Error("api_service.ExpireData FAIL")
    }
}