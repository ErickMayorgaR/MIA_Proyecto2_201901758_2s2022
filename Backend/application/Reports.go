package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

func createPortTd(content string, port string) string {
	return "<td port=\"" + port + "\">" + content + "</td>"
}

type RepCmd struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func (cmd *RepCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "name" {
			cmd.Name = parameter.StringValue
		} else if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "id" {
			cmd.Id = parameter.StringValue
		} else if parameter.Name == "ruta" {
			cmd.Ruta = parameter.StringValue
		}
	}
}

func (cmd *RepCmd) Rep() {
	if cmd.Name == "" || cmd.Id == "" || cmd.Path == "" {
		fmt.Println("Error: faltan parametros obligatorios en el comando rep")
		return
	}

	// CREA LAS CARPETAS PADRE
	parent_path := cmd.Path[:strings.LastIndex(cmd.Path, "/")]

	if err := os.MkdirAll(parent_path, 0777); err != nil {
		log.Fatal(err)
	}

	// OBTENGO EL NOMBRE DEL REPORTE A GENERAR
	report_name := cmd.Path[:strings.LastIndex(cmd.Path, ".")+1]

	// OBTENGO LA PARTICION MONTADA
	partition_m := utils.GlobalList.GetElement(cmd.Id)
	if partition_m.Path == "" {
		fmt.Println("Error: No se puede generar reporte con la particion " + cmd.Id + " ya que no se encuentra montada en RAM")
		return
	}
	// ABRO EL ARCHIVO
	file, err := os.OpenFile(partition_m.Path, os.O_RDWR, 0777)
	// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
	if err != nil {
		log.Fatal("Error ", err)
		return
	}

	// LEO EL MBR
	mbr := utils.MBR{}
	mbr = utils.ReadFileMbr(file, 0)

	// LEO TODO LO RELACIONADO AL SISTEMA DE ARCHIVOS
	super_bloque := utils.SuperBloque{}
	super_bloque = utils.ReadSuperBlock(file, partition_m.Start)

	// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
	var bitinodes = make([]byte, utils.ByteToInt(super_bloque.Inodes_count[:]))
	var bitblocks = make([]byte, utils.ByteToInt(super_bloque.Blocks_count[:]))
	bitinodes = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
	bitblocks = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

	if cmd.Name == "disk" {
		var porcentage float64 = 0

		dotContent := `digraph html { abc [shape=none, margin=0, label=< 
			<TABLE BORDER="1" COLOR="#10a20e" CELLBORDER="1" CELLSPACING="3" CELLPADDING="4">`

		logicas := "\n<TR>"
		all_partitions := "\n<TR>\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\">MBR</TD>\n"

		avalaible_space := 0

		for i := 0; i < 4; i++ {
			if utils.ByteToInt(mbr.Partitions[i].Part_size[:]) != -1 {
				if utils.ByteToString(mbr.Partitions[i].Part_type[:]) != "p" {

					colspan := 2
					temp := utils.EBR{} // GUARDA EL TEMPORAL PARA RECORRER LA LISTA

					// LEO LA PRIMERA PARTICION LOGICA A DONDE APUNTA LA EXTENDIDA Y ASIGNO A TEMP
					temp = utils.ReadEbr(file, utils.ByteToInt(mbr.Partitions[i].Part_start[:]))
					// GRAFICA DE LOGICA
					porcentage = float64(float64((utils.ByteToInt(temp.Part_size[:]))*100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
					logicas += "\n<TD COLOR=\"#87b8a4\">EBR</TD>\n"
					logicas += "\n<TD COLOR=\"#87b8a4\"> Lógica <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

					// MIENTRAS NO LLEGUE AL FINAL DE LA LISTA
					for utils.ByteToInt(temp.Part_next[:]) != -1 {
						colspan += 2
						temp = utils.ReadEbr(file, utils.ByteToInt(temp.Part_next[:]))
						// GRAFICA DE LOGICA
						porcentage = float64((float64(utils.ByteToInt(temp.Part_size[:])) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
						logicas += "\n<TD COLOR=\"#87b8a4\">EBR</TD>\n"
						logicas += "\n<TD COLOR=\"#87b8a4\"> Lógica <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
					}

					// MIENTRAS NO LLEGUE A LA ULTIMA PARTICION
					if i != 3 {
						// GRAFICA DE EXTENDIDA
						porcentage = float64((float64(utils.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
						all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Extendida <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

						if utils.ByteToInt(mbr.Partitions[i+1].Part_size[:]) != -1 {
							// CALCULO ESPACIO VACIO
							avalaible_space = utils.ByteToInt(mbr.Partitions[i+1].Part_start[:]) - (utils.ByteToInt(mbr.Partitions[i].Part_start[:]) + utils.ByteToInt(mbr.Partitions[i].Part_size[:]))
							// CALCULO PORCENTAJE
							porcentage = float64((float64(avalaible_space))*100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:]))
							if porcentage > 0.8 {
								all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
							}
						}
					} else {
						// GRAFICA DE EXTENDIDA
						porcentage = float64((float64(utils.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
						all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Extendida <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
						//all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Extendida <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
						// CALCULO ESPACIO VACIO
						avalaible_space = utils.ByteToInt(mbr.Mbr_size[:]) - (utils.ByteToInt(mbr.Partitions[i].Part_start[:]) + utils.ByteToInt(mbr.Partitions[i].Part_size[:]))
						// CALCULO PORCENTAJE
						porcentage = float64((float64(avalaible_space) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
						if porcentage > 0.8 {
							//all_partitions += "\n<TD COLOR=\"#75e400\" COLSPAN=\"" + strconv.Itoa(colspan) + "\"> Libre <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
							all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
						}
					}
				} else {
					// GRAFICO SOLO LAS PARTICIONES EXISTENTES
					if utils.ByteToInt(mbr.Partitions[i].Part_size[:]) != -1 {
						// MIENTRAS NO LLEGUE A LA ULTIMA PARTICION
						if i != 3 {
							// GRAFICA DE PRIMARIA
							porcentage = float64((float64(utils.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))

							all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Primaria <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"

							if utils.ByteToInt(mbr.Partitions[i+1].Part_size[:]) != -1 {
								// CALCULO ESPACIO VACIO
								avalaible_space = utils.ByteToInt(mbr.Partitions[i+1].Part_start[:]) - (utils.ByteToInt(mbr.Partitions[i].Part_start[:]) + utils.ByteToInt(mbr.Partitions[i].Part_size[:]))
								// CALCULO PORCENTAJE
								porcentage = float64((float64(avalaible_space) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))

								if porcentage > 0.8 {
									all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
								}
							} else {
								// CALCULO ESPACIO VACIO
								avalaible_space = utils.ByteToInt(mbr.Mbr_size[:]) - (utils.ByteToInt(mbr.Partitions[i].Part_start[:]) + utils.ByteToInt(mbr.Partitions[i].Part_size[:]))
								// CALCULO PORCENTAJE
								porcentage = float64((float64(avalaible_space) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))

								if porcentage > 0.8 {
									all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
								}
							}
						} else {
							// GRAFICA DE PRIMARIA
							porcentage = float64((float64(utils.ByteToInt(mbr.Partitions[i].Part_size[:])) * 100.0) / float64(utils.ByteToInt(mbr.Mbr_size[:])))
							all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Primaria <BR/>" + fmt.Sprintf("%.2f", porcentage) + "%</TD>\n"
							// CALCULO ESPACIO VACIO
							avalaible_space = utils.ByteToInt(mbr.Mbr_size[:]) - (utils.ByteToInt(mbr.Partitions[i].Part_start[:]) + utils.ByteToInt(mbr.Partitions[i].Part_size[:]))
							// CALCULO PORCENTAJE
							porcentage = float64((avalaible_space * 100) / utils.ByteToInt(mbr.Mbr_size[:]))

							if porcentage > 0.8 {
								all_partitions += "\n<TD COLOR=\"#75e400\" ROWSPAN=\"3\"> Libre <BR/>" + fmt.Sprint(porcentage) + "%</TD>\n"
							}
						}
					}
				}
			}
		}

		all_partitions += "</TR>\n"
		logicas += "</TR>\n"
		dotContent += all_partitions + logicas + "</TABLE>>];\n}"

		// CREO Y ESCRIBO EL ARCHIVO .dot
		err := ioutil.WriteFile(report_name+"dot", []byte(dotContent), 0777)
		if err != nil {
			log.Fatal(err)
		}

		// GRAFICO EL ARCHIVO .dot CREADO
		utils.GraphDot(report_name+"dot", cmd.Path)
	} else if cmd.Name == "tree" {
		// VARIABLE PARA RECORRER INODOS
		temp_inode := utils.InodeTable{}

		// VARIABLE PARA MOSTRAR TODOS LOS TIPOS DE BLOQUES
		file_block := utils.FileBlock{}
		archive_block := utils.ArchiveBlock{}

		nodes := ""
		blocks := ""
		edges := ""

		dotContent := "digraph {\ngraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\nnode [shape=plain]\nrankdir=LR;"

		// RECORRO INODOS
		for i := 0; i < len(bitinodes); i++ {
			// SI NO ES UN INODO LIBRE
			if bitinodes[i] != '0' {

				// LEO EL INODO
				temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:])+(i*int(unsafe.Sizeof(temp_inode))))
				//fmt.Println(utils.ByteToInt(super_bloque.Inode_start[:]) + (i * int(unsafe.Sizeof(temp_inode))))

				nodes += "inode" + strconv.Itoa(i) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
				nodes += "<tr><td bgcolor=\"#01f5ab\">INODE</td><td bgcolor=\"#01f5ab\">" + strconv.Itoa(i) + "</td></tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("UID", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Uid[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("GID", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Gid[:]), "")
				nodes += "</tr>\n"

				nodes += "<tr>"
				nodes += createPortTd("SIZE", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Size[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("LECTURA", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Atime[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("CREACION", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Ctime[:]), "")
				nodes += "</tr>\n"
				nodes += "<tr>"
				nodes += createPortTd("MODIFICACION", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Mtime[:]), "")
				nodes += "</tr>\n"

				for block_index := 0; block_index < 16; block_index++ {
					nodes += "<tr>"
					nodes += createPortTd("AP"+strconv.Itoa(block_index), "")
					nodes += createPortTd(strconv.Itoa(int(temp_inode.Block[block_index])), "i"+strconv.Itoa(i)+"b"+strconv.Itoa(int(temp_inode.Block[block_index])))
					nodes += "</tr>\n"

					if temp_inode.Block[block_index] != -1 {
						edges += "inode" + strconv.Itoa(i) + ":i" + strconv.Itoa(i) + "b" + strconv.Itoa(int(temp_inode.Block[block_index])) + "->" + "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + ";\n"
						// SI ES UN INODO DE ARCHIVO
						if utils.ByteToString(temp_inode.Type[:]) == "1" {

							archive_block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[block_index])*int(unsafe.Sizeof(archive_block))))

							// ENCABEZADO DEL BLOQUE
							blocks += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
							blocks += "<tr><td bgcolor=\"#f6ec1e\">BLOCK</td><td bgcolor=\"#f6ec1e\">" + strconv.Itoa(int(temp_inode.Block[block_index])) + "</td></tr>\n"
							// ENCABEZADO EL CONTENIDO
							block_content := ""
							for con := 0; con < 64; con++ {
								block_content += string(archive_block.Content[con])
							}
							// LIMPIO MI STRING DE BYTES
							block_content = strings.TrimRight(block_content, "\x00")
							blocks += "<tr><td colspan=\"2\">" + block_content + "</td></tr>\n"
							// CIERRO LA TABLA
							blocks += "</table>>]; \n"
						} else if utils.ByteToString(temp_inode.Type[:]) == "0" {
							file_block = utils.ReadFileBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[block_index])*int(unsafe.Sizeof(file_block))))
							// ENCABEZADO DEL BLOQUE
							blocks += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + " [label=< \n <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
							blocks += "<tr><td bgcolor=\"#f61e73\">BLOCK</td><td bgcolor=\"#f61e73\">" + strconv.Itoa(int(temp_inode.Block[block_index])) + "</td></tr>\n"

							for content := 0; content < 4; content++ {
								blocks += "<tr>" + createPortTd(utils.ByteToString(file_block.Content[content].Name[:]), "") + createPortTd(strconv.Itoa(int(file_block.Content[content].Inodo)), "b"+strconv.Itoa(int(temp_inode.Block[block_index]))+"i"+strconv.Itoa(int(file_block.Content[content].Inodo))) + "</tr>\n"
								if file_block.Content[content].Inodo != -1 && utils.ByteToString(file_block.Content[content].Name[:]) != "." && utils.ByteToString(file_block.Content[content].Name[:]) != ".." {
									edges += "block" + strconv.Itoa(int(temp_inode.Block[block_index])) + ":b" + strconv.Itoa(int(temp_inode.Block[block_index])) + "i" + strconv.Itoa(int(file_block.Content[content].Inodo)) + "->" + "inode" + strconv.Itoa(int(file_block.Content[content].Inodo)) + ";\n"
								}
							}
							blocks += "</table>>]; \n"
						}
					}
				}
				nodes += "<tr>"
				nodes += createPortTd("TIPO", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Type[:]), "")
				nodes += "</tr>\n"

				nodes += "<tr>"
				nodes += createPortTd("PERMISOS", "")
				nodes += createPortTd(utils.ByteToString(temp_inode.Perm[:]), "")
				nodes += "</tr>\n"

				nodes += "</table>>]; \n"
			}
		}
		dotContent += nodes
		dotContent += blocks
		dotContent += edges
		dotContent += "\n}"

		//fmt.Println(dotContent)

		// CREO Y ESCRIBO EL ARCHIVO .dot
		err := ioutil.WriteFile(report_name+"dot", []byte(dotContent), 0777)
		if err != nil {
			log.Fatal(err)
		}
		// GRAFICO EL ARCHIVO .dot CREADO
		utils.GraphDot(report_name+"dot", cmd.Path)
	} else if cmd.Name == "file" {
		// GUARDO EL NOMBRE DEL ARCHIVO
		archive_name := cmd.Ruta[strings.LastIndex(cmd.Ruta, "/")+1 : len(cmd.Ruta)]
		dotContent := "digraph L { node [shape=record fontname=Arial];a  [fillcolor=\"#8bff3a\",style=filled,label=\" Nombre del archivo: " + archive_name + "\\n"

		//exec -path=./test.txt
		// OBTENGO EL ULTIMO NODO DE LA RUTA
		last_inode := utils.GetInodeWithPath(cmd.Ruta, partition_m.Path, partition_m.Start)
		// VARIABLE PARA GUARDAR EL INODO DEL ARCHIVO
		archive_inode := utils.InodeTable{}
		// VARIABLE PARA GUARDAR EL BLOQUE TEMPORAL
		temp_block := utils.FileBlock{}
		exist_archive := false
		// RECORRO EL ULTIMO INODO CON SUS BLOQUE HASTA ENCONTRAR EL NOMBRE DEL ARCHIVO
		for i := 0; i < 16; i++ {
			if last_inode.Block[i] != -1 {
				temp_block = utils.ReadFileBlock(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (int(last_inode.Block[i]) * int(unsafe.Sizeof(temp_block)))))
				for block_i := 0; block_i < 4; block_i++ {
					if utils.ByteToString(temp_block.Content[block_i].Name[:]) == archive_name {
						exist_archive = true
						archive_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:])+(int(temp_block.Content[block_i].Inodo)*int(unsafe.Sizeof(archive_inode))))
					}
				}
			}
		}
		archive_content := ""
		archive_block := utils.ArchiveBlock{}
		if exist_archive {
			for block_i := 0; block_i < 16; block_i++ {
				if archive_inode.Block[block_i] != -1 {
					archive_block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(archive_inode.Block[block_i])*int(unsafe.Sizeof(archive_block))))
					for con := 0; con < 64; con++ {
						archive_content += string(archive_block.Content[con])
					}
					// LIMPIO MI STRING DE BYTES
					archive_content = strings.TrimRight(archive_content, "\x00")
					archive_content += "\\n"
				}
			}
			dotContent += archive_content
			dotContent += "\"]}"
			//fmt.Println(dotContent)

			// CREO Y ESCRIBO EL ARCHIVO .dot
			err := ioutil.WriteFile(report_name+"dot", []byte(dotContent), 0777)
			if err != nil {
				log.Fatal(err)
			}
			// GRAFICO EL ARCHIVO .dot CREADO
			utils.GraphDot(report_name+"dot", cmd.Path)
		} else {
			fmt.Println("Error: no se puede generar el reporte del archivo " + cmd.Ruta + " debido a que no existe en el sistema de archivos")
		}
	}
}
