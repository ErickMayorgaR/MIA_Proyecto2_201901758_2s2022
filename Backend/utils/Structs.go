package utils

type Parameter struct {
	Name        string
	StringValue string
	IntValue    int
}

type Command struct {
	Parameters []Parameter
	Name       string
}

type VoidSpace struct {
	Start int
	Size  int
}

type Partition struct {
	Part_status [3]byte
	Part_type   [4]byte
	Part_fit    [4]byte
	Part_start  [16]byte
	Part_size   [16]byte
	Part_name   [16]byte
}

type MBR struct {
	Mbr_size           [16]byte
	Mbr_fecha_creacion [18]byte
	Mbr_dsk_signature  [16]byte
	Dsk_fit            [4]byte
	Partitions         [4]Partition
}

type EBR struct {
	Part_status [3]byte
	Part_fit    [4]byte
	Part_start  [16]byte
	Part_size   [16]byte
	Part_next   [16]byte
	Part_name   [16]byte
}

// File system part
type SuperBloque struct {
	Filesystem_type   [32]byte
	Inodes_count      [32]byte
	Blocks_count      [32]byte
	Free_blocks_count [32]byte
	Free_inodes_count [32]byte
	Mtime             [23]byte
	Mnt_count         [32]byte
	Magic             [32]byte
	Inode_size        [32]byte
	Block_size        [32]byte
	First_inode       [32]byte
	First_block       [32]byte
	Bm_inode_start    [32]byte
	Bm_block_start    [32]byte
	Inode_start       [32]byte
	Block_start       [32]byte
}

type InodeTable struct {
	Uid   [32]byte
	Gid   [32]byte
	Size  [32]byte
	Atime [23]byte
	Ctime [23]byte
	Mtime [23]byte
	Block [16]int32
	Type  [4]byte
	Perm  [32]byte
}

type ContentFile struct {
	Name  [12]byte
	Inodo int32
}

type FileBlock struct {
	Content [4]ContentFile
}

type ArchiveBlock struct {
	Content [64]byte
}

type User struct {
	Uid      string
	Type     string
	Group    string
	User     string
	Password string
}

type Group struct {
	Gid   string
	Type  string
	Group string
}
