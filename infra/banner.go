package infra

import "fmt"

const (
	banner0 = `
  _____             _    
 |  __ \           | |   
 | |__) | ___  ___ | | __
 |  _  / / _ \/ __|| |/ /
 | | \ \|  __/\__ \|   < 
 |_|  \_\\___||___/|_|\_\
`
	banner1 =`

██████╗ ███████╗███████╗██╗  ██╗
██╔══██╗██╔════╝██╔════╝██║ ██╔╝
██████╔╝█████╗  ███████╗█████╔╝ 
██╔══██╗██╔══╝  ╚════██║██╔═██╗ 
██║  ██║███████╗███████║██║  ██╗
╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝
`
	banner2 = `

 ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄    ▄ 
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌
▐░█▀▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌ ▐░▌ 
▐░▌       ▐░▌▐░▌          ▐░▌          ▐░▌▐░▌  
▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌░▌   
▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░▌    
▐░█▀▀▀▀█░█▀▀ ▐░█▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀█░▌▐░▌░▌   
▐░▌     ▐░▌  ▐░▌                    ▐░▌▐░▌▐░▌  
▐░▌      ▐░▌ ▐░█▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄█░▌▐░▌ ▐░▌ 
▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌
 ▀
`
)

func init() {
	fmt.Println(banner0)
}
