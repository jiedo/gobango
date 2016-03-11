package main

import (
    "fmt"
    "chess"
    "math/rand"
    "os"
)


func main() {
    // 默认实现成回调strategy()模式,
    // 但可以实现成更复杂模式, 符合bot通信协议即可
    sleep_time = 0.1
    show_verbose = false
    chess.G_debug_info = false
    if len(os.Args) >= 2 {
        if "-w" in os.Args {
            // 如果选择白方, 则通知对方START.
            chess.Chess_operate("START")
        }
        if "-v" in os.Args {
            show_verbose = true
        }
        if "-d" in os.Args {
            chess.G_debug_info = true
        }
    }

    while true {
        // 首先读取对方的落子位置, 并写入棋盘
        while bot.Side_this_turn == bot.Your_side {
            h, w = bot.Get_point_of_chessman(bot.Your_side)

        // 检测对方是否获胜
        if bot.Is_winner(bot.Your_side) {
            bot.Light_on_win_points()
            time.sleep(0.1)
            bot.Board_dumps()
            bot.Notes_dumps()
            break
        }
        // 回调自己的策略
        h, w = strategy(bot)
        // 写入棋盘并通知对方
        bot.Put_chessman_at_point(bot.My_side, h, w)
        if show_verbose {
            time.sleep(sleep_time/10)
            bot.Board_dumps()
            time.sleep(sleep_time)
        }
        // 检测自己是否获胜
        if bot.Is_winner(bot.My_side) {
            chess.Chess_log(fmt.Sprintf("%s Win.", chess.ID_TO_NOTE[bot.My_side]))
            bot.Notes_dumps()
            break
        }
}

func strategy(self) {
    // 测试AI
    if self.My_side == chess.BLACK_ID {
        // return strategy4(self, 0, true)
        return strategy6(self, 0, true,
                         max_level_good = 2,
                         max_level_bad = 3)
    } else {
        // return strategy4(self, 0, true)
        return strategy6(self, 0, true,
            max_level_good = 4,
            max_level_bad = 4)
    }
}

func strategy6(self, defence_level, is_dup_enforce,
              max_level_good = 1,
              max_level_bad = 2) {
    // 测试AI
    // 同4,
    // 搜索max_level_good步必胜, 和避免max_level_bad步必败
    // 搜索更多步比较耗时
    //
    // is_dup_enforce: 连珠对附近空白是否有加分
    // defence_level: 防御权重, 越大越重视防御
    //
    // 统计双方所有棋子米字形线条交汇计数最高的空白
    // max(points_score) = max(max(your's + defence),  max(mine))
    //
    all_my_blank_points_count_pair = self.Get_score_of_blanks_for_side(self.My_side,
                                                                       is_dup_enforce)
    for pt, count in all_my_blank_points_count_pair {
        if count > 7 {
            if self.win_test(pt, self.My_side) {
                return pt
            }
        }
    }

    all_your_blank_points_count_pair = self.Get_score_of_blanks_for_side(self.Your_side,
                                                                         is_dup_enforce)
    for pt, count in all_your_blank_points_count_pair {
        if count > 7 {
            if self.win_test(pt, self.Your_side) {
                return pt
            }
        }
    }

    all_blank_points_count = {}
    for pt, count in all_your_blank_points_count_pair {
        if count > 1 {
            all_blank_points_count[pt] = count + defence_level
        }
    }
    for pt, count in all_my_blank_points_count_pair {
        if count > 1 {
            all_blank_points_count[pt] = max(all_blank_points_count.get(pt, 0), count)
        }
    }
    all_blank_points_count_pair = all_blank_points_count.items()
    all_blank_points_count_pair.sort(key=lambda x:x[1])
    all_blank_points_count_pair.reverse()

    // make a batter choice
    for pt, count in all_blank_points_count_pair {
        if self.win_test(pt, self.My_side) {
            return pt
        }
    }
    // make a batter choice
    for level_good in range(1, max_level_good) {
        for pt, count in all_blank_points_count_pair {
            if self.Is_a_good_choice(pt, self.My_side, self.Your_side, level_good) {
                chess.Chess_log(fmt.Sprintf("%s GOOD at: %s", (chess.ID_TO_NOTE[self.My_side],
                    chess.Get_label_of_point(pt))))
                return pt
            }
        }
    }

    blank_points_not_bad = []
    max_deep_bad_point_pt = (0, 0)
    max_deep_bad_point_count = 0
    max_deep_bad_point_level = 0
    // don't make a bad choice
    for pt, count in all_blank_points_count_pair {
        is_bad = false
        for level_bad in range(1, max_level_bad) {
            if self.Is_a_bad_choice(pt, self.My_side, self.Your_side, level_bad) {
                chess.Chess_log(fmt.Sprintf("%s BAD at: %s", (chess.ID_TO_NOTE[self.My_side],
                    chess.Get_label_of_point(pt))))
                if level_bad > max_deep_bad_point_level {
                    max_deep_bad_point_level = level_bad
                    max_deep_bad_point_pt = pt
                    max_deep_bad_point_count = count
                } else if level_bad == max_deep_bad_point_level {
                    if max_deep_bad_point_count < count {
                        max_deep_bad_point_pt = pt
                        max_deep_bad_point_count = count
                    }
                }
                is_bad = true
                break
            }
        }
        if not is_bad {
            blank_points_not_bad += [(pt, count)]
        }
    }
    if blank_points_not_bad {
        _, max_count = blank_points_not_bad[0]
        chess.Chess_log(fmt.Sprintf("points not bad: %d, max_count: %d", (len(blank_points_not_bad), max_count)))
        candidates = [pt for pt, count in blank_points_not_bad if count == max_count]
        pt = candidates[rand.Intn(len(candidates))]
        chess.Chess_log(fmt.Sprintf("%s No.6 give: %s", (chess.ID_TO_NOTE[self.My_side],
            chess.Get_label_of_point(pt))))
        return pt
    }

    chess.Chess_log("no good choice.")
    if all_blank_points_count_pair {
        return max_deep_bad_point_pt
    }
    chess.Chess_log("first point.")
    return (chess.HEIGHT/2, chess.WIDTH/2)
}
