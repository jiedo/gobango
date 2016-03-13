package main


import (
    "fmt"
    "chess"
    "chessbot"	
	"time"
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
    // 默认实现成回调strategy()模式,
    // 但可以实现成更复杂模式, 符合bot通信协议即可
    sleep_time := 10
    show_verbose := false
    chess.G_debug_info = false
    if len(os.Args) >= 2 {
        if string_in("-w", os.Args) {
            // 如果选择白方, 则通知对方START.
            chess.Chess_operate("START")
        }
        if string_in("-v", os.Args) {
            show_verbose = true
        }
        if string_in("-d", os.Args) {
            chess.G_debug_info = true
        }
    }

	bot := chess.Bot{}
    bot.Init_data()
    chess.Chess_log("init blank score.", "DEBUG")
    bot.Get_score_of_blanks_for_side(chess.BLACK_ID, true)
    bot.Get_score_of_blanks_for_side(chess.WHITE_ID, true)
    chess.Chess_log("init ok.", "DEBUG")
	
    for {
        // 首先读取对方的落子位置, 并写入棋盘
        for ;bot.Side_this_turn == bot.Your_side; {
            _, err := bot.Get_point_of_chessman(bot.Your_side)
			if err != nil {
				chess.Chess_log(fmt.Sprintf("error: %s", err), "INFO")				
				continue
			}
		}
        // 检测对方是否获胜
        if bot.Is_winner(bot.Your_side) {
            bot.Light_on_win_points()
            time.Sleep(time.Second)
            bot.Board_dumps()
            bot.Notes_dumps()
            break
        }
        // 回调自己的策略
        my_pt := chessbot.Strategy(&bot)
        // 写入棋盘并通知对方
        bot.Put_chessman_at_point(bot.My_side, my_pt)
        if show_verbose {
            time.Sleep(time.Duration(sleep_time/10) * time.Millisecond)
            bot.Board_dumps()
            time.Sleep(time.Duration(sleep_time) * time.Millisecond)
        }
        // 检测自己是否获胜
        if bot.Is_winner(bot.My_side) {
            chess.Chess_log(fmt.Sprintf("%s Win.", chess.ID_TO_NOTE[bot.My_side]), "INFO")
            bot.Notes_dumps()
            break
        }
	}
}
