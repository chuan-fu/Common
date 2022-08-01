package colorx

import "fmt"

var (
	BlackArrow     = fmt.Sprintf(`%c[1;30m➜ %c[0m`, 0x1B, 0x1B) // 黑色
	RedArrow       = fmt.Sprintf(`%c[1;31m➜ %c[0m`, 0x1B, 0x1B) // 红色
	GreenArrow     = fmt.Sprintf(`%c[1;32m➜ %c[0m`, 0x1B, 0x1B) // 绿色
	YellowArrow    = fmt.Sprintf(`%c[1;33m➜ %c[0m`, 0x1B, 0x1B) // 黄色
	BlueArrow      = fmt.Sprintf(`%c[1;34m➜ %c[0m`, 0x1B, 0x1B) // 蓝色
	PurpleRedArrow = fmt.Sprintf(`%c[1;35m➜ %c[0m`, 0x1B, 0x1B) // 紫红色
	GreenBlueArrow = fmt.Sprintf(`%c[1;36m➜ %c[0m`, 0x1B, 0x1B) // 青蓝色
	WhiteArrow     = fmt.Sprintf(`%c[1;37m➜ %c[0m`, 0x1B, 0x1B) // 白色

	DefaultPrefix   = "-> # "
	RedBluePrefix   = fmt.Sprintf(`%c[1;31m➜ %c[0m%c[1;34m## %c[0m`, 0x1B, 0x1B, 0x1B, 0x1B)
	GreenBluePrefix = fmt.Sprintf(`%c[1;32m➜ %c[0m%c[1;34m## %c[0m`, 0x1B, 0x1B, 0x1B, 0x1B)
)
