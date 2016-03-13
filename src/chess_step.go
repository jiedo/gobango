package main


import (
	"chess"
	"chessbot"
	"os"
)

func string_in(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func main() {
    // 默认实现成回调strtaegy()模式,
    // 但可以实现成更复杂模式, 符合bot通信协议即可
    chess.G_debug_info = false
    if len(os.Args) >= 2 {
        if string_in("-d", os.Args) {
            chess.G_debug_info = true
        }
    }

    board_block := ""
	bot := chess.Bot{}
	
    board_block = `
   - - - - - - - - - - - - - - -
15|. . . . . . . . . . . . . . .|
14|. . . . . . . . . . . . . . .|
13|. . . . . . . . . . . . . . .|
12|. . . O * * * O . . . . . . .|
11|. . . . . * * . . . . . . . .|
10|. . . . . . O . . . . . . . .|
 9|. . . . . . . . . . . . . . .|
 8|. . . . . . . O . . . . . . .|
 7|. . . . . . . . . . . . . . .|
 6|. . . . . . . . . . . . . . .|
 5|. . . . . . . . . . . . . . .|
 4|. . . . . . . . . . . . . . .|
 3|. . . . . . . . . . . . . . .|
 2|. . . . . . . . . . . . . . .|
 1|. . . . . . . . . . . . . . .|
   - - - - - - - - - - - - - - -
   A B C D E F G H I J K L M N O
`
	
    board_block = `
   - - - - - - - - - - - - - - -
15|. . . . . . . . . . . . . . .|
14|. . . . . . . . . . . . . . .|
13|. . . . . . . . . . . . . . .|
12|. . . O * * * O . . . . . . .|
11|. . . . . * * . . . . . . . .|
10|. . . . . . O . . . . . . . .|
 9|. . . . . . . . . . . . . . .|
 8|. . . . . . . O . . . . . . .|
 7|. . . . . . . . . . . . . . .|
 6|. . . . . . . . . . . . . . .|
 5|. . . . . . . . . . . . . . .|
 4|. . . . . . . . . . . . . . .|
 3|. . . . . . . . . . . . . . .|
 2|. . . . . . . . . . . . . . .|
 1|. . . . . . . . . . . . . . .|
   - - - - - - - - - - - - - - -
   A B C D E F G H I J K L M N O
`


    bot.Board_loads(board_block)
    bot.Board_dumps()

	max_level_good := 5
	max_level_bad := 5
    pt := chessbot.Strategy6(&bot, 0, true,
		max_level_good,
		max_level_bad)

    bot.Put_chessman_at_point(bot.Side_this_turn, pt)
    bot.Board_dumps()
}
