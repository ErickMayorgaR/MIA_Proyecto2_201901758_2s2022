package application

import (
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/disk"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/filesystem"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/filesystemadmin"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/partitions"
	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

type ExecCmd struct {
	Path string
}

func (cmd *ExecCmd) AssignParameters(command utils.Command) {
	for _, parameter := range command.Parameters {
		if parameter.Name == "path" {
			cmd.Path = parameter.StringValue
		}
	}
}

func (cmd *ExecCmd) Exec() {
	// obtengo el arbol de comandos
	tree := AnalyzerF(cmd.Path, true)

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
			rep := RepCmd{}
			rep.AssignParameters(element)
			rep.Rep()
		} else if element.Name == "comment" {
			comment := Comment{}
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
			Pause("Pause: Presiona cualquier letra para continuar[*]")
		}
	}
}
