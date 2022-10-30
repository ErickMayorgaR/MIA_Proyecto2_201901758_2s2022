package utils

type MyUser struct {
	Logged       int
	Uid          string
	User_name    string
	Pwd          string
	Grp          string
	Id_partition string
	Gid          string
}

var GlobalUser = MyUser{}
