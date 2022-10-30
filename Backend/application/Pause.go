package application

import (
	"fmt"
	"time"
)

// FUNCION PARA DETENER PROCESOS
func counter() {
	i := 0
	for {
		time.Sleep(time.Second * 1)
		i++
	}
}
func Pause(message string) {
	//var v string
	//fmt.Scanln(&v)

	go counter()
	fmt.Println(message)
	fmt.Scanln()

	/*ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)
	fmt.Println("Pause: Presiona cualquier letra para continuar[*]")
	for {
		stdin, _ := <-ch
		// AL DETECTAR CUALQUIER ENTRADA DE TECLADO ROMPE EL CICLO
		if stdin != "" {
			fmt.Println("\nSaliendo del Pause...")
			stdin = ""
			break
		}

	}*/
}
