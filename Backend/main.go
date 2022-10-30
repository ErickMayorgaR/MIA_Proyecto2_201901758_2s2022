package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/application"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/disk"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/filesystem"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/filesystemadmin"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/partitions"
)

type event struct {
	ID          string `json:"ID"`
	Tittle      string `json:"Titulo"`
	Description string `json:"Data"`
}

var contenido = []event{
	{ID: "1", Tittle: "Primer Item", Description: "Respuesta"},
	{ID: "2", Tittle: "Segundo Item", Description: "Respuesta"},
	{ID: "3", Tittle: "Tercer Item", Description: "Respuesta"},
}

func main() {
	startProcess()

	/*
		server := startServer(":9090")

		err := server.ListenAndServe()

		if err != nil {
			panic(err)
		}
	*/

}

func startServer(addr string) *http.Server {
	initRoutes()
	return &http.Server{
		Addr: addr,
	}

}

func initRoutes() {
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			giveInformation(w, r)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			giveInformation(w, r)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	})

}

func giveInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	json.NewEncoder(w).Encode(contenido)

}

func addInformation(w http.ResponseWriter, r *http.Request) {

}

/*
	router := gin.Default()
	router.GET("/get", giveInformation)
	router.Run("localhost:9090")


	func giveInformation(context *gin.Context) {
	fmt.Println("Informacion enviada")
	fmt.Println(contenido)

	context.IndentedJSON(http.StatusOK, contenido)

}

*/

func startProcess() {
	// Dirección máquina virtual /usr/local/go/src/github.com/PR2_MIA
	for true {
		CallClear()
		fmt.Println("------------------------------")
		fmt.Println("------------------------------")
		// LEO ENTRADA CON ESPACIOS HASTA ENCONTRAR UN SALTO DE LINEA
		in := bufio.NewReader(os.Stdin)
		inputCommand, _ := in.ReadString('\n')

		// OBTENGO EL ARBOL DE COMANDOS
		tree := application.AnalyzerF(inputCommand, false)
		// BORRO LOS COMANDOS INNECESARIOS
		tree = tree[:2]

		//exec -path=./test.txt
		//exec -path=./test1.txt

		for _, element := range tree {
			if element.Name == "mkdisk" {
				mkdisk := disk.MkdiskCmd{}
				mkdisk.AssignParameters(element)
				mkdisk.Mkdisk()
			} else if element.Name == "rmdisk" {
				rmdisk := disk.RmdiskCmd{}
				rmdisk.AssignParameters(element)
				rmdisk.Rmdisk()
			} else if element.Name == "fdisk" {
				fdisk := partitions.FdiskCmd{}
				fdisk.AssignParameters(element)
				fdisk.Fdisk()
			} else if element.Name == "mount" {
				mount := partitions.MountCmd{}
				mount.AssignParameters(element)
				mount.Mount()
			} else if element.Name == "mkfs" {
				mkfs := filesystem.MkfsCmd{}
				mkfs.AssignParameters(element)
				mkfs.Mkfs()
			} else if element.Name == "mkdir" {
				mkdir := filesystemadmin.MkdirCmd{}
				mkdir.AssignParameters(element)
				mkdir.Mkdir()
			} else if element.Name == "mkfile" {
				mkfile := filesystemadmin.MkfileCmd{}
				mkfile.AssignParameters(element)
				mkfile.Mkfile()
			} else if element.Name == "rep" {
				rep := application.RepCmd{}
				rep.AssignParameters(element)
				rep.Rep()
			} else if element.Name == "comment" {
				comment := application.Comment{}
				comment.AssignParameters(element)
				comment.ShowComment()
			} else if element.Name == "login" {
				login := filesystem.LoginCmd{}
				login.AssignParameters(element)
				login.Login()
			} else if element.Name == "logout" {
				filesystem.Logout()
			} else if element.Name == "mkgrp" {
				mkgrp := filesystem.MkgrpCmd{}
				mkgrp.AssignParameters(element)
				mkgrp.Mkgrp()
			} else if element.Name == "rmgrp" {
				rmgrp := filesystem.RmgrpCmd{}
				rmgrp.AssignParameters(element)
				rmgrp.Rmgrp()
			} else if element.Name == "mkuser" {
				mkuser := filesystem.MkuserCmd{}
				mkuser.AssignParameters(element)
				mkuser.Mkuser()
			} else if element.Name == "rmusr" {
				rmuser := filesystem.RmusrCmd{}
				rmuser.AssignParameters(element)
				rmuser.Rmusr()
			} else if element.Name == "pause" {
				application.Pause("Pause: Presiona cualquier letra para continuar[*]")
			} else if element.Name == "exec" {
				exec := application.ExecCmd{}
				exec.AssignParameters(element)
				exec.Exec()
			}
		}
		application.Pause("Fin de ejecución del script")
	}
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
