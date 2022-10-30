package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func GetLetter(number int) string {
	letter := ""
	for ch := 97; ch <= 122; ch++ {
		if ch == (number + 96) {
			letter = string(rune(ch))
		}
	}
	return letter
}

type PartitionMounted struct {
	DiskName      string // DISCO DONDE SE ENCUENTRA LA PARTICION MONTADA
	PartitionName string // NOMBRE DE LA PARTICION MONTADA
	Id            string // ID DE LA PARTICION MONTADA
	Start         int    // INICIO DE LA PARTICION MONTADA EN EL DISCO
	Status        int    // ATRIBUTO PARA SABER SI ESTA MONTADA LA PARTICION
	Type          string // TIPO DE PARTICION MONTADA
	Fit           string // TIPO DE FIT DE LA PARTICION MONTADA
	Size          int    // TAMANIO DE LA PARTICION MONTADA
	Path          string // RUTA DONDE SE ENCUENTRA EL DISCO DE LA PARTICION
}

type Element struct {
	DiskName          string
	DiskNumber        int
	Path_disk         string // RUTA DEL DISCO
	PartitionsMounted []PartitionMounted
}

type ListPartitionsMounted struct {
	Partitions []Element
}

// FUNCION PARA CREAR NUEVO ELEMENTO
func NewElement(disk_name string, disk_number int, path_disk string) Element {
	slice := make([]PartitionMounted, 1)
	return Element{disk_name, disk_number, path_disk, slice}
}

// FUNCION PARA CREAR UNA NUEVA LISTA DE PARTICIONES
func NewRam() ListPartitionsMounted {
	slice := make([]Element, 1)
	return ListPartitionsMounted{slice}
}

func (list *ListPartitionsMounted) MountPartition(disk_name string, partition_name string, start int, path string, type_partition string, fit string, size int) {
	idNumber := -1
	idLeter := ""

	// BUSCA SI YA SE HA MONTADO ALGUNA PARTICION DEL DISCO
	for i, partition := range list.Partitions {
		if partition.DiskName == disk_name {
			for _, par := range list.Partitions[i].PartitionsMounted {
				if par.PartitionName == partition_name && par.Status == 1 {
					fmt.Println("Error: la particion " + partition_name + " ya está montada en ram " + disk_name)
					return
				}
			}
			idNumber = partition.DiskNumber
			idLeter = GetLetter(len(partition.PartitionsMounted) + 1)
			id := "57" + strconv.Itoa(idNumber) + idLeter
			newPartition := PartitionMounted{disk_name, partition_name, id, start, 1, type_partition, fit, size, path}
			list.Partitions[i].PartitionsMounted = append(list.Partitions[i].PartitionsMounted, newPartition)
			break
		}
	}

	// SI NO SE ENCONTRO EL DISCO LO CREO
	if idNumber == -1 {
		idNumber = len(list.Partitions)
		idLeter = "a"
		id := "57" + strconv.Itoa(idNumber) + idLeter
		element := NewElement(disk_name, len(list.Partitions), path)
		list.Partitions = append(list.Partitions, element)
		newPartition := PartitionMounted{disk_name, partition_name, id, start, 1, type_partition, fit, size, path}
		list.Partitions[len(list.Partitions)-1].PartitionsMounted = append(list.Partitions[len(list.Partitions)-1].PartitionsMounted, newPartition)
	}

	//fmt.Println(idNumber, idLeter, disk_name)
}

func (list *ListPartitionsMounted) GetElement(myId string) PartitionMounted {
	res := PartitionMounted{}
	// BUSCA SI YA SE HA MONTADO ALGUNA PARTICION DEL DISCO
	for _, partition := range list.Partitions {
		for _, par := range partition.PartitionsMounted { //list.Partitions[i].PartitionsMounted {
			if par.Id == myId {
				res = par
				break
			}
		}
	}
	return res
}

var GlobalList ListPartitionsMounted = NewRam()

// FUNCION PARA LEER BITMPAS DEL DISCO
func ReadBitMap(file *os.File, position int, bitmap_size int) []byte {
	var bitmap = make([]byte, bitmap_size)

	size := int(bitmap_size)
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &bitmap)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return bitmap
}

// FUNCION PARA LEER SUPERBLOQUE
func ReadSuperBlock(file *os.File, position int) SuperBloque {
	var super_bloque = SuperBloque{}

	size := int(unsafe.Sizeof(super_bloque))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &super_bloque)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return super_bloque
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadInode(file *os.File, position int) InodeTable {
	var inode = InodeTable{}

	size := int(unsafe.Sizeof(inode))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &inode)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return inode
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadFileBlock(file *os.File, position int) FileBlock {
	var block = FileBlock{}

	size := int(unsafe.Sizeof(block))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &block)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return block
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadArchiveBlock(file *os.File, position int) ArchiveBlock {
	var block = ArchiveBlock{}

	size := int(unsafe.Sizeof(block))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &block)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return block
}

// FUNCION PARA LEER INODOS DEL DISCO EN LA PARTICION
func ReadFileMbr(file *os.File, position int) MBR {
	var mbr = MBR{}

	size := int(unsafe.Sizeof(mbr))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &mbr)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return mbr
}

// FUNCION PARA LEER PARTICIONES DEL DISCO
func ReadPartition(file *os.File, position int) Partition {
	var partition = Partition{}

	size := int(unsafe.Sizeof(partition))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &partition)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return partition
}

// FUNCION PARA LEER EBR DEL DISCO
func ReadEbr(file *os.File, position int) EBR {
	var ebr = EBR{}

	size := int(unsafe.Sizeof(ebr))
	file.Seek(int64(position), 0)
	data := ReadBytes(file, size)
	buffer := bytes.NewBuffer(data)
	err1 := binary.Read(buffer, binary.BigEndian, &ebr)
	if err1 != nil {
		log.Fatal("Error ", err1)
	}
	return ebr
}

// FUNCION PARA BUSCAR EL ULTIMO INODO DE UNA RUTA
func GetInodeWithPath(path string, real_path string, start int) InodeTable {
	//fmt.Println("entraaa en el file")
	if GlobalUser.Logged == -1 {
		fmt.Println("Error: Para utilizar mkfile necesitas estar logueado")
		return InodeTable{}
	}

	// ABRO EL ARCHIVO
	file, err := os.OpenFile(real_path, os.O_RDWR, 0777)
	// VERIFICACION DE ERROR AL ABRIR EL ARCHIVO
	if err != nil {
		log.Fatal("Error ", err)
		return InodeTable{}
	}

	// OBTENGO TODAS LAS CARPETAS PADRES ANTES DEL ARCHIVO
	parent_path := path[:strings.LastIndex(path, "/")]

	//current_date := GetDate()

	// OBTENGO TODAS LAS CARPETAS PADRES DEL ARCHIVO
	routes := strings.Split(parent_path, "/")

	//saltar_busqueda := false
	exist_route := false

	// SI EL TAMANIO DE LAS RUTAS SEPARADAS POR / ES
	// CERO QUIERE DECIR QUE EL ARCHIVO SE DEBE CREAR EN LA RAIZ
	if len(routes) == 0 && path == "/" {
		exist_route = true
	} else {
		temp := []string{"/"}
		temp = append(temp, routes[1:]...)
		routes = temp
	}

	// LEO EL SUPERBLOQUE
	super_bloque := SuperBloque{}
	super_bloque = ReadSuperBlock(file, start)

	// CREACION DE ARRAY PARA ALMACENAR LOS BITMPAS
	var bitinodes = make([]byte, ByteToInt(super_bloque.Inodes_count[:]))
	var bitblocks = make([]byte, ByteToInt(super_bloque.Blocks_count[:]))
	bitinodes = ReadBitMap(file, ByteToInt(super_bloque.Bm_inode_start[:]), len(bitinodes))
	bitblocks = ReadBitMap(file, ByteToInt(super_bloque.Bm_block_start[:]), len(bitblocks))

	// VERIFICA QUE EXISTAN LAS CARPETAS ANTES DEL ARCHIVO

	temp_inode := InodeTable{}

	//LEO EL PRIMER INODO
	temp_inode = ReadInode(file, ByteToInt(super_bloque.Inode_start[:]))

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
			if temp_inode.Block[pointerIndex] != -1 && ByteToString(temp_inode.Type[:]) == "0" {
				file_block := FileBlock{}
				file_block = ReadFileBlock(file, (ByteToInt(super_bloque.Block_start[:]) + (int(temp_inode.Block[pointerIndex]) * int(unsafe.Sizeof(file_block)))))
				// RECORRE LOS PUNTEROS DE LOS BLOQUES
				for blockIndex := 0; blockIndex < 4; blockIndex++ {
					if file_block.Content[blockIndex].Inodo != -1 {
						if ByteToString(file_block.Content[blockIndex].Name[:]) == routes[path_index] {
							// ELIMINO LAS RUTAS QUE YA ESTAN CREADAS PARA QUE QUEDEN SOLO LAS RESTANTES
							if len(remaining_routes) == len(routes) {
								remaining_routes = RemoveIndex(remaining_routes, 0)
								remaining_routes = RemoveIndex(remaining_routes, 0)
							} else {
								remaining_routes = RemoveIndex(remaining_routes, 0)
							}
							temp_inode = ReadInode(file, ByteToInt(super_bloque.Inode_start[:])+(int(file_block.Content[blockIndex].Inodo)*int(unsafe.Sizeof(temp_inode))))
							exist_route = true
							exist_path = true
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
		temp_inode = ReadInode(file, ByteToInt(super_bloque.Inode_start[:]))
		exist_route = true
	}
	//VERIFICACION DE EXISTENCIA DE RUTAS
	if !exist_route {
		fmt.Println("No existe la ruta indicada")
		return InodeTable{}
	}
	return temp_inode
}

// FUNCION PARA ESCRIBIR BLOQUES DE CARPETAS
func WriteFileBlocks(file *os.File, position int, fileBlock FileBlock) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &fileBlock)
	WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR BLOQUES DE ARCHIVOS
func WriteArchiveBlocks(file *os.File, position int, archiveBlock ArchiveBlock) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &archiveBlock)
	WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR TABLAS DE INODOS
func WriteInodes(file *os.File, position int, inode InodeTable) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &inode)
	WriteBytes(file, bufferControlBlocks.Bytes())

}

// FUNCION PARA ESCRIBIR EL SUPERBLOQUE EN EL ARCHIVO
func WriteSuperBlock(file *os.File, position int, super_bloque SuperBloque) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &super_bloque)
	WriteBytes(file, bufferControlBlocks.Bytes())
}

// FUNCION PARA ESCRIBIR BITMAPS
func WriteBitmap(file *os.File, position int, bitmap []byte) {
	file.Seek(int64(position), 0)
	var bufferControlBlocks bytes.Buffer
	binary.Write(&bufferControlBlocks, binary.BigEndian, &bitmap)
	WriteBytes(file, bufferControlBlocks.Bytes())
}
