package ucmevents

import (
	"fmt"
)

type UcmEvento struct{
	Modulo string
	Mensaje string 
	Parametros  map[string]string 
}

type UcmPublicador struct {	
	
}

// Este método debe sustituirse por uno adecuado pra publicar eventos en la UCM
func (UcmPublicador) Publish(e UcmEvento) error{
	
	return fmt.Errorf("No se pudo publicar evento")
}