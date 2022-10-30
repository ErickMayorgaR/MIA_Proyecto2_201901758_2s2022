package filesystemadmin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type MkdirCmd struct {
	Path string
	P    string
}

func (cmd *MkdirCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		} else if parameter.Name == "-p" {
			cmd.P = parameter.StringValue
		}
	}
}

func (cmd *MkdirCmd) Mkdir() {
	if utils.GlobalUser.Logged == -1 {
		fmt.Println("Error: para utilizar mkdir necesitas estar logueado")
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

	// LEO EL SUPERBLOQUE DE LA PARTICION
	super_bloque := utils.SuperBloque{}
	size_superbloque := int(unsafe.Sizeof(super_bloque))
	file.Seek(int64(partition_m.Start), 0)
	data_superbloque := utils.ReadBytes(file, size_superbloque)
	buffer_superbloque := bytes.NewBuffer(data_superbloque)
	err1 := binary.Read(buffer_superbloque, binary.BigEndian, &super_bloque)
	if err1 != nil {
		log.Fatal("Error ", err1)
		return
	}

	temp := strings.Split(cmd.Path, "/")
	routes := []string{"/"}
	if cmd.Path != "/" {
		routes = append(routes, temp[1:]...)
	}

	// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
	var bitinodes = make([]byte, utils.ByteToInt(super_bloque.Inodes_count[:]))
	var bitblocks = make([]byte, utils.ByteToInt(super_bloque.Blocks_count[:]))
	bitinodes = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
	bitblocks = utils.ReadBitMap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

	// CREA LA RAIZ
	if len(routes) == 1 {
		// CREO EL INODO RAIZ Y LLENO SUS ATRIBUTOS
		root_inode := utils.InodeTable{}
		copy(root_inode.Uid[:], []byte(utils.GlobalUser.Uid))
		copy(root_inode.Gid[:], []byte(utils.GlobalUser.Gid))
		copy(root_inode.Size[:], []byte(strconv.Itoa(0)))
		copy(root_inode.Atime[:], []byte(utils.GetDate()))
		copy(root_inode.Ctime[:], []byte(utils.GetDate()))
		copy(root_inode.Mtime[:], []byte(utils.GetDate()))
		// INICIALIZO LOS APUNTADORES
		for i := 0; i < 16; i++ {
			root_inode.Block[i] = -1
		}
		copy(root_inode.Type[:], []byte("0"))
		copy(root_inode.Perm[:], []byte(strconv.Itoa(664)))

		// APUNTO HACIA EL PRIMER BLOQUE DE CARPETA
		root_inode.Block[0] = 0

		// CREO EL PRIMER BLOQUE DE CARPETA
		block := utils.FileBlock{}
		for s := 0; s < 4; s++ {
			block.Content[s].Inodo = -1
			copy(block.Content[s].Name[:], []byte(""))
		}
		copy(block.Content[0].Name[:], []byte("."))
		copy(block.Content[1].Name[:], []byte(".."))
		block.Content[0].Inodo = 0
		block.Content[1].Inodo = 0

		// MODIFICO LOS BITMAPS
		bitinodes[0] = '1'
		bitblocks[0] = '1'

		// MODIFICO ATRIBUTOS DEL SUPERBLOQUE
		copy(super_bloque.First_inode[:], []byte(strconv.Itoa(1)))
		copy(super_bloque.First_block[:], []byte(strconv.Itoa(1)))
		copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.First_block[:])-1)))
		copy(super_bloque.Free_blocks_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.First_block[:])-1)))

		// ESCRIBO EL BLOQUE DE CARPETA
		utils.WriteFileBlocks(file, utils.ByteToInt(super_bloque.Block_start[:]), block)
		// ESCRIBO EL INODO PRINCIPAL
		utils.WriteInodes(file, utils.ByteToInt(super_bloque.Inode_start[:]), root_inode)

		// REESCRIBO SUPERBLOQUE
		utils.WriteSuperBlock(file, partition_m.Start, super_bloque)
		// REESCRIBO BITMAPS
		utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), bitinodes)
		utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)
	} else {
		temp_inode := utils.InodeTable{}
		var index_temp_inode int
		catch_inode_index := 0
		// LEO EL PRIMER INODO
		temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:]))
		exist_route := false

		// VECTOR PARA GUARDAR LAS RUTAS QUE FALTAN POR CREARSE
		var remaining_routes = make([]string, len(routes))
		// CREO UNA COPIA PARA QUE NO SE ALTERE EL ROUTE
		copy(remaining_routes, routes)

		var ver = make([]int, 1)
		ver[0] = 0

		// RECORRE LA RUTA
		for path_index := 0; path_index < len(routes); path_index++ {
			exist_path := false
			// RECORRE PUNTEROS DEL INODO
			for pointerIndex := 0; pointerIndex < 16; pointerIndex++ {

				// RECORRO SOLO LOS BLOQUES DE LOS INODOS DE TIPO CARPETA
				if temp_inode.Block[pointerIndex] != -1 && utils.ByteToString(temp_inode.Type[:]) == "0" {
					file_block := utils.FileBlock{}
					index_temp_inode = int(temp_inode.Block[pointerIndex]) // GUARDO EL INIDICE DEL INDICE TEMP
					file_block = utils.ReadFileBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[pointerIndex])*int(unsafe.Sizeof(file_block))))
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
								// LEO EL INODO TEMPORAL
								temp_inode = utils.ReadInode(file, utils.ByteToInt(super_bloque.Inode_start[:])+(int(file_block.Content[blockIndex].Inodo)*int(unsafe.Sizeof(temp_inode))))
								//fmt.Println(utils.ByteToInt(super_bloque.Inode_start[:]) + (int(file_block.Content[blockIndex].Inodo) * int(unsafe.Sizeof(temp_inode))))
								//index_temp_inode = int(file_block.Content[blockIndex].Inodo)
								ver = append(ver, int(file_block.Content[blockIndex].Inodo))

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
		if !exist_route {
			// IGUALO EL index_temp_inode AL INDICE DE LA ULTIMA RUTA ENCONTRADA
			index_temp_inode = catch_inode_index
			if cmd.P != "" || len(remaining_routes) == 1 {

				// SI SOLAMENTE ESTA CREADA LA RAIZ LA ELIMINO DE LAS RUTAS RESTANTES
				if len(remaining_routes) == len(routes) {
					remaining_routes = utils.RemoveIndex(remaining_routes, 0)
				}

				// CREO LAS RUTAS RESTANTES
				// RECORRO LAS RUTAS RESTANTES
				for route_index := 0; route_index < len(remaining_routes); route_index++ {
					for pointer := 0; pointer < 16; pointer++ {
						has_space := false
						var indice_encontrado int

						actual_block := utils.FileBlock{}

						// VERIFICA QUE SEA UN PUNTERO SIN UTILIZAR
						if temp_inode.Block[pointer] != -1 {
							actual_block = utils.ReadFileBlock(file, utils.ByteToInt(super_bloque.Block_start[:])+(int(temp_inode.Block[pointer])*int(unsafe.Sizeof(actual_block))))

							for block_index := 0; block_index < 4; block_index++ {
								// VALIDACION DE PUNTERO LIBRE
								if actual_block.Content[block_index].Inodo == -1 { // RECORRO EL BLOQUE LIBRE
									indice_encontrado = block_index // OBTENGO EL INDICE DISPONIBLE DEL BLOQUE
									has_space = true
									break
								}
								has_space = false
							}
						} else {
							// SI ES -1 ENTONCES LLEGO AL PRIMER PUNTERO DISPONIBLE
							// ENTONCES SE CREA UN NEVO BLOQUE Y SE APUNTA AL PUNTERO
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
							if pointer == 0 {
								real_block.Content[0].Inodo = int32(index_temp_inode)
								copy(real_block.Content[0].Name[:], []byte("."))
								real_block.Content[1].Inodo = int32(ver[len(ver)-2])
								copy(real_block.Content[1].Name[:], []byte(".."))
								block_pointer = 2
							}
							// APUNTO EL INODO ACTUAL AL BLOQUE
							temp_inode.Block[pointer] = int32(free_block)

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
							copy(newInode.Type[:], []byte("0"))
							copy(newInode.Perm[:], []byte(strconv.Itoa(664)))
							// MODFICO LOS BITMAPS Y CONTADORES
							bitinodes[free_inode] = '1'
							bitblocks[free_block] = '1'
							copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))
							copy(super_bloque.Free_blocks_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_blocks_count[:])-1)))
							// APUNTO EL PRIMER PUNTERO DEL BLOQUE AL INODO CREADO
							// Y GUARDO EL NOMBRE DEL NUEVO DIRECTORIO CREADO
							real_block.Content[block_pointer].Inodo = int32(free_inode)

							copy(real_block.Content[block_pointer].Name[:], []byte(remaining_routes[route_index]))

							// ESCRIBO EL INODO NUEVO
							utils.WriteInodes(file, utils.ByteToInt(super_bloque.Inode_start[:])+(free_inode*int(unsafe.Sizeof(newInode))), newInode)
							// AQUIII ESTA EL ERROR
							utils.WriteInodes(file, utils.ByteToInt(super_bloque.Inode_start[:])+(index_temp_inode*int(unsafe.Sizeof(temp_inode))), temp_inode)

							// ESCRIBO EL BLOQUE
							utils.WriteFileBlocks(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (free_block * int(unsafe.Sizeof(real_block)))), real_block)
							//fmt.Println((utils.ByteToInt(super_bloque.Block_start[:]) + (free_block * int(unsafe.Sizeof(real_block)))))
							// REESCRIBO EL SUPERBLOQUE
							utils.WriteSuperBlock(file, partition_m.Start, super_bloque)

							// REESCRIBO BITMAPS
							utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), bitinodes)
							utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)

							/*for i := 0; i < 4; i++ {
								fmt.Println(real_block.Content[i].Inodo, utils.ByteToString(real_block.Content[i].Name[:]))
							}*/
							// RECORRO EL ARBOL
							temp_inode = newInode
							index_temp_inode = free_inode
							// VEEEER
							ver = append(ver, free_inode)
							break
						}

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
							copy(newInode.Type[:], []byte("0"))
							copy(newInode.Perm[:], []byte("664"))
							// MODIFICO LOS BITMPAS
							bitinodes[free_inode] = '1'
							// MODIFICO ATRIBUTOS DEL SUPERBLOQUE
							copy(super_bloque.Free_inodes_count[:], []byte(strconv.Itoa(utils.ByteToInt(super_bloque.Free_inodes_count[:])-1)))

							// APUNTO EL APUNTADOR DEL BLOQUE DISPONIBLE AL NUEVO INODO
							// Y GUARDO EL NOMBRE DEL DIRECTORIO ACTUAL
							actual_block.Content[indice_encontrado].Inodo = int32(free_inode)
							copy(actual_block.Content[indice_encontrado].Name[:], []byte(remaining_routes[route_index]))

							// ESCRIBO EL INODO
							utils.WriteInodes(file, (utils.ByteToInt(super_bloque.Inode_start[:]) + (free_inode * int(unsafe.Sizeof(newInode)))), newInode)
							// ESCRIBO EL BLOQUE
							utils.WriteFileBlocks(file, (utils.ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointer]) * int(unsafe.Sizeof(actual_block)))), actual_block)
							//fmt.Println(utils.ByteToInt(super_bloque.Block_start[:]))

							// REESCRIBO EL SUPERBLOQUE
							utils.WriteSuperBlock(file, partition_m.Start, super_bloque)

							// REESCRIBO BITMAPS
							utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_inode_start[:]), bitinodes)
							utils.WriteBitmap(file, utils.ByteToInt(super_bloque.Bm_block_start[:]), bitblocks)
							// RECORRO EL ARBOL
							temp_inode = newInode
							index_temp_inode = free_inode
							// VEEEER
							ver = append(ver, free_inode)
							break
						} else {
							continue
						}
					}
				}
			} else {
				fmt.Println("Error: la ruta no existe, para crearla usa el parametro -p")
				return
			}
		}
	}

	file.Close()
}
