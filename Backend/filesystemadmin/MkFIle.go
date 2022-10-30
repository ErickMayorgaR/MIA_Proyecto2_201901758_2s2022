package filesystemadmin

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type MkfileCmd struct {
	Path    string
	R       string
	Size    int
	Cont    string
	AnyText string
}

func (cmd *MkfileCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "-r" {
			cmd.R = parameter.StringValue
		} else if parameter.Name == "size" {
			cmd.Size = parameter.IntValue
		} else if parameter.Name == "cont" {
			cmd.Cont = parameter.StringValue
		}
	}
}

func (cmd *MkfileCmd) Mkfile() {
	//fmt.Println("entraaa en el file")
	if utils.GlobalUser.Logged == -1 {
		fmt.Println("Error: Para utilizar mkfile necesitas estar logueado")
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

	if cmd.Path != "" {
		// GUARDO EL NOMBRE DEL ARCHIVO
		archive_name := cmd.Path[strings.LastIndex(cmd.Path, "/")+1 : len(cmd.Path)]

		// OBTENGO TODAS LAS CARPETAS PADRES ANTES DEL ARCHIVO
		parent_path := cmd.Path[:strings.LastIndex(cmd.Path, "/")]

		//current_date := utils.GetDate()

		// OBTENGO TODAS LAS CARPETAS PADRES DEL ARCHIVO
		routes := strings.Split(parent_path, "/")

		//saltar_busqueda := false
		exist_route := false

		// SI EL TAMANIO DE LAS RUTAS SEPARADAS POR / ES
		// CERO QUIERE DECIR QUE EL ARCHIVO SE DEBE CREAR EN LA RAIZ
		if len(routes) == 0 && cmd.Path == "/" {
			exist_route = true
		} else {
			temp := []string{"/"}
			temp = append(temp, routes[1:]...)
			routes = temp
		}

		// LEO EL SUPERBLOQUE
		super_bloque := utils.SuperBloque{}
		super_bloque = utils.ReadSuperBlock(file, partition_m.Start)

		// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
		var bitinodes = make([]byte, utils.ByteToInt(super_bloque.Inodes_count[:]))
		var bitblocks = make([]byte, utils.ByteToInt(super_bloque.Blocks_count[:]))
		bitinodes = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
		bitblocks = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

		// VERIFICA QUE EXISTAN LAS CARPETAS ANTES DEL ARCHIVO

		temp_inode := utils.InodeTable{}
		var index_temp_inode int
		catch_inode_index := 0
		//LEO EL PRIMER INODO
		temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:]))

		// VECTOR PARA GUARDAR LAS RUTAS QUE FALTAN POR CREARSE
		var remaining_routes = make([]string, len(routes))
		// CREO UNA COPIA PARA QUE NO SE ALTERE EL ROUTE
		copy(remaining_routes, routes)

		// RECORRE LA RUTA
		for path_index := 0; path_index < len(routes); path_index++ {
			exist_path := false
			// RECORRE LOS PUNTEROS DEL INODO
			for pointerIndex := 0; pointerIndex < 16; pointerIndex++ {
				// RECORRO SOLO LOS BLOQUES DE LOS INODOS DE TIPO CARPETA
				if temp_inode.Block[pointerIndex] != -1 && utils.ByteToString(temp_inode.Type[:]) == "0" {
					file_block := utils.FileBlock{}
					index_temp_inode = int(temp_inode.Block[pointerIndex]) // GUARDO EL INDICE DEL TEMP
					file_block = utils.ReadFileBlock(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointerIndex]) * int(unsafe.Sizeof(file_block)))))
					// RECORRE LOS PUNTEROS DE LOS BLOQUES
					for blockIndex := 0; blockIndex < 4; blockIndex++ {
						if file_block.Content[blockIndex].Inodo != -1 {
							if utils.ByteToString(file_block.Content[blockIndex].Name[:]) == routes[path_index] {
								// ELIMINO LAS RUTAS QUE YA ESTAN CREADAS PARA QUE QUEDEN SOLO LAS RESTANTES
								if len(remaining_routes) == len(routes) {
									remaining_routes = utils.RemoveIndex(remaining_routes, 0)
									remaining_routes = utils.RemoveIndex(remaining_routes, 0)
								} else {
									remaining_routes = utils.RemoveIndex(remaining_routes, 0)
								}
								temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:])+(int(file_block.Content[blockIndex].Inodo)*int(unsafe.Sizeof(temp_inode))))
								index_temp_inode = int(file_block.Content[blockIndex].Inodo)
								exist_route = true
								exist_path = true
								// ATRAPA EL ULTIMO NUMBERO DE INODO QUE SE CREO
								catch_inode_index = int(file_block.Content[blockIndex].Inodo)
							}
						}
					}
				}
			}
			if !exist_path {
				exist_route = false
			}
		}
		// VALIDACION PARA SABER SI EL ARCHIVO SE CREA EN LA RAIZ
		if routes[0] == "/" && len(routes) == 1 {
			temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:]))
			index_temp_inode = 0
			exist_route = true
		}
		//VERIFICACION DE EXISTENCIA DE RUTAS
		if !exist_route {
			// IGUALO EL index_temp_inode AL INDICE DE LA ULTIMA RUTA ENCONTRADA
			index_temp_inode = catch_inode_index
			//fmt.Println("Entraa", remaining_routes)
			// CREA LAS RUTAS FALTANTES
			if cmd.R != "" {
				// CREA LAS RUTAS
				c := MkdirCmd{}
				c.Path = parent_path
				c.P = "-p"
				c.Mkdir()
				// VUELVE A EJECUTAR EL MKFILE DESPUES DE CREAR LAS RUTAS RESTANTES
				d := MkfileCmd{}
				d.AnyText = cmd.AnyText
				d.Cont = cmd.Cont
				d.Path = cmd.Path
				d.R = cmd.R
				d.Size = cmd.Size
				d.Mkfile()
				return
			} else {
				fmt.Println("Error: la ruta " + cmd.Path + " en mkfile no existe intenta utilizando el parametro -r")
				file.Close()
				return
			}
		} else {
			// BUSCO EL PUNTERO LIBRE DEL INODO
			for pointer_index := 0; pointer_index < 15; pointer_index++ {
				has_space := false
				var indice_encontrado int // GUARDA EL INDICE DEL BLOQUE QUE ESTE DISPONIBLE
				actual_block := utils.FileBlock{}
				// VERIFICA QUE SEA UN PUNTERO SIN UTILIZAR
				if temp_inode.Block[pointer_index] != -1 {
					actual_block = utils.ReadFileBlock(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointer_index]) * int(unsafe.Sizeof(actual_block)))))
					for block_index := 0; block_index < 4; block_index++ { // RECORRO EL BLOQUE
						// VALIDO QUE EXISTE UN PUNTERO LIBRE
						if actual_block.Content[block_index].Inodo == -1 {
							indice_encontrado = block_index // OBTENGO EL INDICE DISPONIBLE DEL BLOQUE
							has_space = true
							break
						}
						has_space = false
					}
				} else {
					// SI ES -1 ENTONCES LLEGO AL PRIMER PUNTERO DISPONIBLE
					// ENTONCES SE CREA UN NUEVO BLOQUE Y SE APUNTA AL PUNTERO
					// ACTUAL EL CUAL APUNTA A UN NUEVO INODO
					var free_block int
					var free_inode int
					real_block := utils.FileBlock{}
					// LEO BITMAP DE BLOQUES
					bitblocks = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))
					// LEO BITMAP DE INODOS
					bitinodes = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
					// CALCULO LA POSICION DEL BLOQUE
					for bit_index := 0; bit_index < len(bitblocks); bit_index++ {
						if bitblocks[bit_index] == '0' {
							free_block = bit_index
							break
						}
					}
					// INICIALIZO LOS PUNTEROS Y NOMBRES DEL BLOQUE
					for i_block := 0; i_block < 4; i_block++ {
						real_block.Content[i_block].Inodo = -1
						copy(real_block.Content[i_block].Name[:], []byte(""))
					}

					block_pointer := 0
					// SI ES EL PRIMER APUNTADOR DEL INODO LOS PRIMEROS DOS REGISTROS APUNTAN AL PADRE
					if pointer_index == 0 {
						real_block.Content[0].Inodo = int32(index_temp_inode)
						copy(real_block.Content[0].Name[:], []byte("."))
						real_block.Content[1].Inodo = int32(index_temp_inode)
						copy(real_block.Content[1].Name[:], []byte(".."))
						block_pointer = 2
					}
					// APUNTO EL INODO ACTUAL AL BLOQUE
					temp_inode.Block[pointer_index] = int32(free_block)

					// CALCULO LA POSICION DEL INODO LIBRE
					for bit_index := 0; bit_index < len(bitinodes); bit_index++ {
						if bitinodes[bit_index] == '0' {
							free_inode = bit_index
							break
						}
					}
					// CREO EL INODO
					newInode := utils.InodeTable{}
					copy(newInode.Uid[:], []byte(utils.GlobalUser.Uid))
					copy(newInode.Gid[:], []byte(utils.GlobalUser.Gid))
					copy(newInode.Size[:], []byte(strconv.Itoa(0)))
					copy(newInode.Atime[:], []byte(utils.GetDate()))
					copy(newInode.Ctime[:], []byte(utils.GetDate()))
					copy(newInode.Mtime[:], []byte(utils.GetDate()))
					// INICIALIZO LOS APUNTADORES
					for i := 0; i < 16; i++ {
						newInode.Block[i] = -1
					}

					copy(newInode.Type[:], []byte("1"))
					copy(newInode.Perm[:], []byte(strconv.Itoa(664)))
					// MODIFICO LOS BITMAPS Y CONTADORES
					bitinodes[free_inode] = '1'
					bitblocks[free_block] = '1'
					copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))
					copy(super_bloque.Free_blocks_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_blocks_count[:])-1)))
					// APUNTO EL PRIMER PUNTERO DEL BLOQUE AL INODO CREADO
					// Y GUARDO EL NOMBRE DEL NUEVO DIRECTORIO CREADO
					real_block.Content[block_pointer].Inodo = int32(free_inode)
					copy(real_block.Content[block_pointer].Name[:], []byte(archive_name))

					// ESCRIBO EL INODO NUEVO
					utils.WriteInodes(file, (utils.ByteToInt(super_bloque.Inode_start[:]) + (free_inode * int(unsafe.Sizeof(newInode)))), newInode)
					utils.WriteInodes(file, (utils.ByteToInt(super_bloque.Inode_start[:]) + (index_temp_inode * int(unsafe.Sizeof(temp_inode)))), temp_inode)
					// ESCRIBO EL BLOQUE
					utils.WriteFileBlocks(file, utils.ByteToInt(super_bloque.Block_start[:])+(free_block*int(unsafe.Sizeof(real_block))), real_block)

					// REESCRIBO EL SUPERBLOQUE
					utils.WriteSuperBlock(file, partition_m.Start, super_bloque)

					// REESCRIBO BITMAPS
					utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), bitinodes)
					utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)

					// RECORRO EL ARBOL
					temp_inode = newInode
					index_temp_inode = free_inode
					break
				}
				// VERIFICA QUE EXISTA ESPACIO Y DETIENE LA BUSQUEDA DE BLOQUES
				if has_space {
					var free_inode int
					// BUSCA EL INODO LIBRE EN EL BITMAP DE INODOS
					bitinodes = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
					for inode_i := 0; inode_i < len(bitinodes); inode_i++ {
						if bitinodes[inode_i] == '0' {
							free_inode = inode_i
							break
						}
					}
					// CREO EL INODO
					newInode := utils.InodeTable{}
					copy(newInode.Uid[:], []byte(utils.GlobalUser.Uid))
					copy(newInode.Gid[:], []byte(utils.GlobalUser.Gid))
					copy(newInode.Size[:], []byte(strconv.Itoa(0)))
					copy(newInode.Atime[:], []byte(utils.GetDate()))
					copy(newInode.Ctime[:], []byte(utils.GetDate()))
					copy(newInode.Mtime[:], []byte(utils.GetDate()))
					// INICIALIZO LOS APUNTADORES
					for i := 0; i < 16; i++ {
						newInode.Block[i] = -1
					}
					copy(newInode.Type[:], []byte("1"))
					copy(newInode.Perm[:], []byte(strconv.Itoa(664)))
					// MODIFICO LOS BITMAPS
					bitinodes[free_inode] = '1'
					// MODIFICO LOS ATRIBUTOS DEL SUPERBLOQUE
					copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))
					// APUNTO EL APUNTADOR DEL BLOQUE DISPONIBLE AL NEVO INODO
					// Y GUARDO EL NOMBRE DEL DIRECTORIO ACTUAL
					actual_block.Content[indice_encontrado].Inodo = int32(free_inode)
					copy(actual_block.Content[indice_encontrado].Name[:], []byte(archive_name))

					// ESCRIBO EL INODO
					utils.WriteInodes(file, (utils.ByteToInt(super_bloque.Inode_start[:]) + (free_inode * int(unsafe.Sizeof(newInode)))), newInode)
					// ESCRIBO EL BLOQUE
					utils.WriteFileBlocks(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointer_index]) * int(unsafe.Sizeof(actual_block)))), actual_block)

					// REESCRIBO EL SUPERBLOQUE
					utils.WriteSuperBlock(file, partition_m.Start, super_bloque)

					// REESCRIBO BITMAPS
					utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), bitinodes)
					utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)

					// RECORRO EL ARBOL
					temp_inode = newInode
					index_temp_inode = free_inode
					break
				} else {
					continue
				}
			}
		}

		content_size := 0 // VARIABLE PARA SABER EL TAMANIO DEL CONTENIDO
		content := ""     // VARIABLE PARA GUARDAR EL CONTENIDO

		// SI HAY CONTENIDO CALCULO SU TAMANIO
		if cmd.Cont != "" {
			nombreArchivo := cmd.Cont
			content = utils.ReadFile(nombreArchivo)
			content_size = len(content)
		}

		// SI LE PASO TEXTO POR A PARTE HACE EL CALCULO
		if cmd.AnyText != "" {
			content = cmd.AnyText
			content_size = len(cmd.AnyText)
		}

		// VALIDACION DE TAMANIO NEGATIVO EN PARAMETRO SIZE
		if cmd.Size < 0 {
			fmt.Println("Error: el parametro size no puede ser negativo")
			file.Close()
			return
		}

		// ESCRIBO EL ARCHIVO SEGUN EL TAMANIO INDICADO
		if cmd.Size >= 0 || content_size >= 0 {

			if cmd.Size >= content_size {
				copy(temp_inode.Size[:], []byte(strconv.Itoa(cmd.Size)))
			} else {
				copy(temp_inode.Size[:], []byte(strconv.Itoa(content_size)))
			}

			var number_blocks int
			if utils.ByteToInt(temp_inode.Size[:]) > 64 {
				number_blocks = (utils.ByteToInt(temp_inode.Size[:]) / 64)
			} else {
				number_blocks = 1
			}

			if utils.ByteToInt(temp_inode.Size[:]) == 0 {
				number_blocks = 0
			}

			if ((utils.ByteToInt(temp_inode.Size[:]) % 64) != 0) && (utils.ByteToInt(temp_inode.Size[:]) > 64) {
				number_blocks++
			}

			for block_index := 0; block_index < number_blocks; block_index++ {
				var free_block_index int
				// BUSCO EL BLOQUE LIBRE EN EL BITMAP DE BLOQUES
				for bit := 0; bit < len(bitblocks); bit++ {
					if bitblocks[bit] == '0' {
						free_block_index = bit
						break
					}
				}

				// ESCRIBO EN EL BLOQUE DE ARCHIVO LOS CARACTERES DEL 1 AL 9
				// Y EL CONTENIDO SI ES QUE TIENE EN CONT O ANYTEXT
				bloqueArchivo := utils.ArchiveBlock{}
				temp_inode.Block[block_index] = int32(free_block_index)
				number := 0
				for caracter_index := 0; caracter_index < 64; caracter_index++ {
					// SI NO LE ENVIO EL VALOR INTERNAMENTE ENTONCES NO LLENA EL BLOQUE CON NUMEROS DEL 1 AL 9
					if cmd.AnyText == "" {
						bloqueArchivo.Content[caracter_index] = byte(number + 48)
						// SI ES IGUAL A 9 VUELVE A 0 Y SI NO AUMENTA
						if number == 9 {
							number = 0
						} else {
							number++
						}
					}

					if len(content) >= 1 /*0*/ {
						if len(content) > 1 {
							bloqueArchivo.Content[caracter_index] = content[0]
							// BORRO EL PRIMER CARACTER DEL CONTENIDO HASTA VACIARLA
							content = content[1:]
						}
					}
				}

				// MODIFICO ATRIBUTOS DEL SUPERBLOQUE
				copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))
				bitblocks[free_block_index] = '1'

				// REESCRIBO EL SUPERBLOQUE EN LA PARTICION
				utils.WriteSuperBlock(file, partition_m.Start, super_bloque)

				// REESCRIBO EL INODO TEMPORAL QUE ES EL DEL ARCHIVO
				utils.WriteInodes(file, (utils.ByteToInt(super_bloque.Inode_start[:]) + (index_temp_inode * int(unsafe.Sizeof(temp_inode)))), temp_inode)

				// ESCRIBO EL BLOQUE DE ARCHIVO EN EL DISCO
				utils.WriteArchiveBlocks(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (free_block_index * int(unsafe.Sizeof(bloqueArchivo)))), bloqueArchivo)
				//ver := utils.ArchiveBlock{}
				//ver = utils.ReadArchiveBlock(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (free_block_index * int(unsafe.Sizeof(bloqueArchivo)))))
				//fmt.Println(ver.Content, (utils.ByteToInt(super_bloque.Block_start[:]) + (free_block_index * int(unsafe.Sizeof(bloqueArchivo)))))
				// REESCRIBO EL BITMAP DE BLOQUES
				utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)
			}

		}
		file.Close()

	} else {
		fmt.Println("Error: el parametro path es obligatorio en el comando mkfile")
	}
}
