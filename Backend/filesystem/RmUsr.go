package filesystem

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type RmusrCmd struct {
	Usuario string
}

func (cmd *RmusrCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "usuario" {
			cmd.Usuario = parameter.StringValue
		}
	}
}

func (cmd *RmusrCmd) Rmusr() {

	if cmd.Usuario != "" {

		// VALIDA QUE EXISTA UN USUARIO LOGUEADO
		if utils.GlobalUser.Logged == -1 {
			fmt.Println("Error: Para crear un grupo necesitas estar logueado")
			return
		} else if utils.GlobalUser.User_name != "root" {
			fmt.Println("Error: Para crear un grupo necesitas estar logueado con el usuario root")
			return
		}

		// VARIABLE CON TODA LA INFORMACION DE LA PARTICION MONTADA
		partition_m := utils.GlobalList.GetElement(utils.GlobalUser.Id_partition)

		// ABRO EL ARCHIVO
		file, err := os.OpenFile(partition_m.Path, os.O_RDWR, 0777)
		// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
		if err != nil {
			log.Fatal("Error ", err)
			return
		}

		// LEO EL SUPERBLOQUE
		super_bloque := utils.SuperBloque{}
		super_bloque = utils.ReadSuperBlock(file, partition_m.Start)

		// LEO EL PRIMER INODO QUE ES EL QUE CONTIENE EL ARCHIVO DE USUARIOS
		users_inode := utils.InodeTable{}
		users_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:])+int(unsafe.Sizeof(users_inode)))

		archive_block := utils.ArchiveBlock{}

		users_archive_content := ""

		//actual_block_index := 0

		for block_i := 0; block_i < 16; block_i++ {
			if users_inode.Block[block_i] != -1 {
				archive_block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_i])*int(unsafe.Sizeof(archive_block))))
				// CONCATENO QUITANDO EL SALTO DE LINEA DERECHO PARA QUE NO DE ERROR
				users_archive_content += strings.TrimRight(utils.ByteToString(archive_block.Content[:]), "")
				//actual_block_index = block_i
			}
		}

		// ALMACENO TODOS LOS GRUPOS Y USUARIOS SEPARADOS POR UN SALTO
		all := strings.Split(users_archive_content, "\n")

		// ARREGLOS PARA GUARDAR LOS GRUPOS Y USUARIOS POR SEPARADO
		var groups = make([]utils.Group, 0)
		var users = make([]utils.User, 0)

		// RECORRO TODOS LOS USUAIROS Y GRUPOSY LOS SEPARO
		for i := 0; i < len(all); i++ {
			if all[i] != "" {
				temp := strings.Split(all[i], ",")
				if temp[1] == "G" {
					groups = append(groups, utils.Group{temp[0], temp[1], temp[2]})
				} else if temp[1] == "U" {
					users = append(users, utils.User{temp[0], temp[1], temp[2], temp[3], temp[4]})
				}
			}
		}

		// STRING PARA GUARDAR EL NUEVO STRING DE GRUPOS Y USUARIOS SIN EL USUSARIO QUE SE BORRO
		new_string := ""

		exist_user_in := false
		for i := 0; i < len(users); i++ {
			// SI ES EL MISMO USUARIO ACTIVA LA BANDERA Y NO CONCATENA
			if cmd.Usuario == users[i].User {
				if users[i].Uid != "0" {
					new_string += "\n" + "0" + "," + users[i].Type + "," + users[i].Group + "," + users[i].User + "," + users[i].Password + "\n"
					exist_user_in = true
				}
			} else {
				new_string += users[i].Uid + "," + users[i].Type + "," + users[i].Group + "," + users[i].User + "," + users[i].Password + "\n"
			}
		}

		for i := 0; i < len(groups); i++ {
			new_string += groups[i].Gid + "," + groups[i].Type + "," + groups[i].Group + "\n"

		}

		// VALIDA QUE EL GRUPO YA ESTE CREADO EN LA PARTICION
		if !exist_user_in {
			fmt.Println("Error: el usuario " + cmd.Usuario + " no puede eliminarse ya que no existe en la particion " + partition_m.PartitionName)
		} else {
			// QUITO DEL STRING TODOS LOS SALTOS DE LINEA A LA DERECHA
			users_archive_content = strings.TrimRight(new_string, "\n")

			// LEO BITMAP DE BLOQUES
			var bitblocks = make([]byte, utils.ByteToInt(super_bloque.Blocks_count[:]))
			bitblocks = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

			caracter_count := 0           // CONTADOR PARA POSICIONARME EN EL STRING
			block_index := 0              // INDICE PARA EL BLOQUE ACTUAL
			block := utils.ArchiveBlock{} // BLOQUE ACTUAL
			block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_index])*int(unsafe.Sizeof(block))))
			// LIMPIO EL BLOQUE DE CONTENIDO
			block.Content = utils.NewByte()

			// REDIMENSIONO EL ARCHIVO DE USUARIOS
			copy(users_inode.Size[:], strconv.Itoa(len(users_archive_content)))

			// RECORRO EL STRING CON LOS GRUPOS Y USARIOS
			for len(users_archive_content) != 0 {
				if caracter_count == 63 {
					// ESCRIBO EL BLOQUE
					utils.WriteArchiveBlocks(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_index])*int(unsafe.Sizeof(block))), block)

					block_index++
					caracter_count = 0
					if int(users_inode.Block[block_index]) != -1 {
						block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_index])*int(unsafe.Sizeof(block))))
						// LIMPIO EL BLOQUE DE CONTENIDO
						block.Content = utils.NewByte()
					} else {
						var free_block_index int
						// BUSCO EL BLOQUE LIBRE EN EL BITMAP DE BLOQUES
						for bit := 0; bit < len(bitblocks); bit++ {
							if bitblocks[bit] == '0' {
								free_block_index = bit
								break
							}
						}
						//block = utils.ReadArchiveBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[free_block_index])*int(unsafe.Sizeof(block))))
						block = utils.ArchiveBlock{}
						users_inode.Block[block_index] = int32(free_block_index)
						// REESCRIBO EL INODO QUE CONTIENE LOS BLOQUES DEL ARCHIVO DE USUARIO
						utils.WriteInodes(file, utils.ByteToInt(super_bloque.Inode_start[:])+int(unsafe.Sizeof(users_inode)), users_inode)
						// MODIFICO ATRIBUTOS DEL SUPERBLOQUE
						copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))
						bitblocks[free_block_index] = '1'
						// REESCRIBO EL SUPERBLOQUE EN LA PARTICION
						utils.WriteSuperBlock(file, partition_m.Start, super_bloque)
						// REESCRIBO EL BITMAP DE BLOQUES
						utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)
					}
				}
				// GUARDO EL CARACTER EN EL CARACTER DEL BLOQUE
				block.Content[caracter_count] = users_archive_content[0]
				users_archive_content = users_archive_content[1:]
				caracter_count++
			}

			// ESCRIBO EL BLOQUE
			utils.WriteArchiveBlocks(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(users_inode.Block[block_index])*int(unsafe.Sizeof(block))), block)
		}

	} else {
		fmt.Println("Error: el parametro id es obligatorio en el comando login")
	}

}
