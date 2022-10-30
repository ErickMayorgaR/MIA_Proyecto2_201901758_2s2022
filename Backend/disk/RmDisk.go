package disk

import (
	"fmt"
	"os"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type RmdiskCmd struct {
	Path string
}

func (cmd *RmdiskCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		}
	}
}

func (cmd *RmdiskCmd) Rmdisk() {
	err := os.Remove(cmd.Path)
	if err != nil {
		fmt.Println("Error: al eliminar el disco en la ruta " + cmd.Path)
	}
}
