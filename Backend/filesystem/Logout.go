package filesystem

import (
	"fmt"

	"github.com/ErickMayorgaR/MIA_Proyecto2_201901758_2s2022/utils"
)

func Logout() {
	if utils.GlobalUser.Logged != -1 {
		// CIERRO SESION
		utils.GlobalUser.Logged = -1
		utils.GlobalUser.Uid = ""
		utils.GlobalUser.User_name = ""
		utils.GlobalUser.Pwd = ""
		utils.GlobalUser.Grp = ""
		utils.GlobalUser.Id_partition = ""
		utils.GlobalUser.Gid = ""
	} else {
		fmt.Println("Error: no se puede realizar el logout ya que no hay ningun usuario logueado actualmente")
	}
}
