package main


import (
	"chess"
	"chessbot"	
)

func main() {
    // 默认实现成回调strtaegy()模式,
    // 但可以实现成更复杂模式, 符合bot通信协议即可
    chess.G_debug_info = true

	bot := chess.Bot{}
	bot.Init_data()
    board_block := `
   - - - - - - - - - - - - - - -
15|. . . . . . . . . . . . . . .|
14|. . . . . . . . . . . . . . .|
13|. . . . . . . * . . . . . . .|
12|. . . . . . O . O . . . . . .|
11|. . * * * * O * . * . . . . .|
10|. . . . * . . . . . . . . . .|
 9|. . . O . O O . . . . . . . .|
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

	max_level_good := 3
	max_level_bad := 3

    pt := chessbot.Strategy6(&bot, 0, true,
                     max_level_good,
                     max_level_bad)

    bot.Put_chessman_at_point(bot.Side_this_turn, pt)
    bot.Board_dumps()
}
