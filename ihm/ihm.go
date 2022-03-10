package ihm

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func Avisar(mensaje string) {
	color.Cyan("[AVISO] ")
	color.HiCyan("%s \n", mensaje)
}

func Problema(mensaje string) {
	color.Red("[ERROR] ")
	color.HiRed("%s \n", mensaje)
}

func Letrero(mensaje string) {
	color.HiBlue("%s \n", mensaje)
}

func Texto(mensaje string) {
	color.White("%s \n", mensaje)
}

func Suceso(mensaje string) {
	color.HiGreen("%s \n", mensaje)
}

func Impacto(mensaje string) {
	color.Set(color.FgWhite, color.Bold)
	fmt.Printf("%s \n", mensaje)
	color.Unset()
}

func Evento(mensaje string) {
	ahora := time.Now()
	color.HiGreen("%s ", ahora.Format("2006.01.02 15:04:05"))
	color.White("%s \n", mensaje)
}
