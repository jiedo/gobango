package chessbot


import (
    "fmt"
    "chess"
    "math/rand"
)



func Strategy(self *chess.Bot) chess.Point{
    // 测试AI

	max_level_good := 3
	max_level_bad := 3		
	
    if self.My_side == chess.BLACK_ID {
        // return strategy4(self, 0, true)
        return Strategy6(self, 0, true,
                         max_level_good,
                         max_level_bad)
    } else {
        // return strategy4(self, 0, true)
		max_level_good = 3
		max_level_bad = 3		
        return Strategy6(self, 0, true,
            max_level_good,
            max_level_bad)
    }
}

func Strategy6(self *chess.Bot, defence_level int, is_dup_enforce bool,
              max_level_good int,
              max_level_bad int) chess.Point{
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
    all_my_blank_points_count_pair := self.Get_score_of_blanks_for_side(self.My_side, is_dup_enforce)
    for _, ppair := range all_my_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value		
		if count > 7 {
            if self.Win_test(pt, self.My_side) {
                return pt
            }
        }
    }

    all_your_blank_points_count_pair := self.Get_score_of_blanks_for_side(self.Your_side, is_dup_enforce)
    for _, ppair := range all_your_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value		
        if count > 7 {
            if self.Win_test(pt, self.Your_side) {
                return pt
            }
        }
    }

    all_blank_points_count := make(map[chess.Point]int)
    for _, ppair := range all_your_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value		
        if count > 1 {
            all_blank_points_count[pt] = count + defence_level
        }
    }
    for _, ppair := range all_my_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value		
        if count > 1 {
			if count_tmp, ok := all_blank_points_count[pt]; ok {
				if count_tmp > count {
					all_blank_points_count[pt] = count_tmp
				}
			} else {
				all_blank_points_count[pt] = count
			}
        }
    }

	all_blank_points_count_pair := chess.Rank_by_point_count(all_blank_points_count)	
    // make a batter choice
    for _, ppair := range all_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value
        if count < 4 {
			continue
		}
        if self.Win_test(pt, self.My_side) {
            return pt
        }
    }
    // make a batter choice
	for level_good:=1; level_good<max_level_good; level_good++ {		
		for _, ppair := range all_blank_points_count_pair {
			pt, count := ppair.Key, ppair.Value
			if count < 1 {
				continue
			}
            if self.Is_a_good_choice(pt, self.My_side, self.Your_side, level_good) {
                chess.Chess_log(fmt.Sprintf("%s GOOD at: %s", chess.ID_TO_NOTE[self.My_side],
                    chess.Get_label_of_point(pt)), "INFO")
                return pt
            }
        }
    }

    blank_points_not_bad := make(map[chess.Point]int)
    max_deep_bad_point_pt := chess.Point{0, 0}
    max_deep_bad_point_count := 0
    max_deep_bad_point_level := 0
    // don't make a bad choice
	for _, ppair := range all_blank_points_count_pair {
		pt, count := ppair.Key, ppair.Value		
        is_bad := false
        for level_bad:=1; level_bad<max_level_bad; level_bad++ {
            if self.Is_a_bad_choice(pt, self.My_side, self.Your_side, level_bad) {
                chess.Chess_log(fmt.Sprintf("%s BAD at: %s", chess.ID_TO_NOTE[self.My_side],
                    chess.Get_label_of_point(pt)), "INFO")
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
        if ! is_bad {
            blank_points_not_bad[pt] = count
        }
    }

	
	if len(blank_points_not_bad) > 0 {
		blank_points_not_bad_pair := chess.Rank_by_point_count(blank_points_not_bad)		
		// to fix get max
		top_point := blank_points_not_bad_pair[0]
        chess.Chess_log(fmt.Sprintf("points not bad: %d, max_count: %d", len(blank_points_not_bad), top_point.Value), "INFO")

		candidates := make([]chess.Point, len(blank_points_not_bad_pair))
		for _, ppair := range blank_points_not_bad_pair {
			pt, count := ppair.Key, ppair.Value
			if count == top_point.Value {
				candidates = append(candidates, pt)
			}
		}
        pt := candidates[rand.Intn(len(candidates))]
        chess.Chess_log(fmt.Sprintf("%s No.6 give: %s", chess.ID_TO_NOTE[self.My_side],
            chess.Get_label_of_point(pt)), "INFO")
        return pt
    }

    chess.Chess_log("no good choice.", "INFO")
    if len(all_blank_points_count_pair) > 0 {
        return max_deep_bad_point_pt
    }
    chess.Chess_log("first point.", "INFO")
    return chess.Point{chess.HEIGHT/2, chess.WIDTH/2}
}
