package disk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type MkdiskCmd struct {
	Size int
	Fit  string
	Unit string
	Path string
}

func (cmd *MkdiskCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "size" {
			cmd.Size = parameter.IntValue
		} else if parameter.Name == "fit" {
			cmd.Fit = parameter.StringValue
		} else if parameter.Name == "unit" {
			cmd.Unit = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		}
	}
}

func (cmd *MkdiskCmd) Mkdisk() {
	if cmd.Size != -1 {
		// CREA LAS CARPETAS PADRE
		parent_path := cmd.Path[:strings.LastIndex(cmd.Path, "/")]
		if err := os.MkdirAll(parent_path, 0700); err != nil {
			log.Fatal(err)
		}

		// CREA EL ARCHIVO
		disk_file, err := os.Create(cmd.Path)
		if err != nil {
			log.Fatal(err)
		}

		// RELLENA EL ARCHIVO CON CEROS
		multiplicator := 1024
		if cmd.Unit == "k" {
			multiplicator = 1
		}
		var temporal int8 = 0
		// RELLENO MI BUFFER CON CEROS
		var binario bytes.Buffer
		for i := 0; i < 1024; i++ {
			binary.Write(&binario, binary.BigEndian, &temporal)
		}

		for i := 0; i < cmd.Size*multiplicator; i++ {
			utils.WriteBytes(disk_file, binario.Bytes())
		}

		// CREO EL MBR
		MBR := utils.MBR{}

		// INICIALIZO LAS PARTICIONES DEL MBR
		for i := 0; i < 4; i++ {
			copy(MBR.Partitions[i].Part_name[:], "")
			copy(MBR.Partitions[i].Part_status[:], "0")
			copy(MBR.Partitions[i].Part_type[:], "P")
			copy(MBR.Partitions[i].Part_start[:], "-1")
			copy(MBR.Partitions[i].Part_size[:], "-1")
			copy(MBR.Partitions[i].Part_fit[:], []byte(cmd.Fit))
		}

		// ASIGNACION DE ATRIBUTOS DEL MBR
		copy(MBR.Dsk_fit[:], []byte(cmd.Fit))
		copy(MBR.Mbr_size[:], []byte(strconv.Itoa(cmd.Size*multiplicator*1024)))
		copy(MBR.Mbr_dsk_signature[:], []byte(strconv.Itoa(utils.GetRandom())))
		copy(MBR.Mbr_fecha_creacion[:], []byte(utils.GetDate()))

		// ESCRIBO EL MBR EN EL DISCO
		disk_file.Seek(0, 0)
		var bufferControl bytes.Buffer
		binary.Write(&bufferControl, binary.BigEndian, &MBR)
		utils.WriteBytes(disk_file, bufferControl.Bytes())

		disk_file.Close()

	} else {
		fmt.Println("Error: el parametro size es obligatorio en mkdisk")
	}
}
