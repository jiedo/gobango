package chess

import (
    "fmt"
    "os"
	"bufio"
    "strings"
    "errors"
    "sort"
    "strconv"
)


const (
    WIDTH int = 15
    HEIGHT = 15
    WIN_NUM = 5
)

const (
    POINT_NOTE int = 0
    POINT_NEED_UPDATE = 12

    LEG_INFO_IDX_DUP_SUM_SCORE = 13
    LEG_INFO_IDX_SUM_SCORE = 14

    LEG_INFO_N = 15
)
type Direction int
const (
    DIRECTION_ROW Direction= 0
    DIRECTION_WEST = 1
    DIRECTION_EAST = 2

    DIRECTION_COL = 3
    DIRECTION_SOUTH = 4
    DIRECTION_NORTH = 5

    DIRECTION_DOWN = 6
    DIRECTION_NORTHWEST = 7
    DIRECTION_SOUTHEAST = 8

    DIRECTION_UP = 9
    DIRECTION_SOUTHWEST = 10
    DIRECTION_NORTHEAST = 11
)

const (
    BLANK string = "."
    WHITE = "O"
    BLACK = "*"
    WHITE_WIN = "\033[32mO\033[0m"
    BLACK_WIN = "\033[32m*\033[0m"

    OP_PUT = "PUT"
    BOARD_MARKS string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    BOARD_MARKS_LENGTH = len(BOARD_MARKS)
)

type GoSide int
const (
    WHITE_ID GoSide = 0
    BLACK_ID = 1
    BLANK_ID = 2
    BLACK_WIN_ID = 3
    WHITE_WIN_ID = 4
)
var (

    MARK_WIN = map[GoSide]GoSide{
        WHITE_ID: WHITE_WIN_ID,
        BLACK_ID: BLACK_WIN_ID,
    }

    OPPOSITE_SIDE = map[GoSide]GoSide {
        WHITE_ID: BLACK_ID,
        BLACK_ID: WHITE_ID,
    }

    NOTE_TO_ID = map[string]GoSide{
        BLANK: BLANK_ID,
        WHITE: WHITE_ID,
        BLACK: BLACK_ID,
    }

    ID_TO_NOTE = []string{WHITE, BLACK, BLANK, BLACK_WIN, WHITE_WIN}

    OPPOSITE_DIRECTION = map[Direction]Direction{
        DIRECTION_UP: DIRECTION_DOWN,
        DIRECTION_DOWN: DIRECTION_UP,
        DIRECTION_ROW: DIRECTION_COL,
        DIRECTION_COL: DIRECTION_ROW,

        DIRECTION_WEST: DIRECTION_EAST,
        DIRECTION_EAST: DIRECTION_WEST,
        DIRECTION_SOUTH: DIRECTION_NORTH,
        DIRECTION_NORTH: DIRECTION_SOUTH,
        DIRECTION_NORTHWEST: DIRECTION_SOUTHEAST,
        DIRECTION_SOUTHEAST: DIRECTION_NORTHWEST,
        DIRECTION_SOUTHWEST: DIRECTION_NORTHEAST,
        DIRECTION_NORTHEAST: DIRECTION_SOUTHWEST,
    }

    G_debug_info bool = false
)
// assert (BOARD_MARKS_LENGTH > WIDTH)
// assert (BOARD_MARKS_LENGTH > HEIGHT)

func idtoa(point_w int) string{
    if point_w < 0 || point_w > BOARD_MARKS_LENGTH {
        return string(BOARD_MARKS[0])
    }
    return string(BOARD_MARKS[point_w])
}


func atoid(mark_w string) int{
    return strings.Index(BOARD_MARKS, string(mark_w))
}


func Get_label_of_point(pt Point) string{
    return fmt.Sprintf("%s%d", idtoa(pt.W), pt.H+1)
}


func Chess_log(msg string, level string) {
    if level == "DEBUG" && !G_debug_info {
        return
    }
    fmt.Fprintln(os.Stderr, msg)
}


func Chess_operate(op string) {
    Chess_log(op, "INFO")
    fmt.Println(op)
}


type Point struct {
    H int
    W int
}

func Rank_by_point_count(points_score map[Point]int) PairList{
  pl := make(PairList, len(points_score))
  i := 0
  for k, v := range points_score {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

type Pair struct {
  Key Point
  Value int
}

type PairList []Pair
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }






type Bot struct {
    started bool
    win_points [HEIGHT]Point
    win_points_num int
    center_point Point
    center_side GoSide
    test_side GoSide

    My_side GoSide
    Your_side GoSide
    Side_this_turn GoSide

    Board [HEIGHT][WIDTH]GoSide

    Board_leg_type [HEIGHT][WIDTH][2][LEG_INFO_N]int
    direction_legtype_count []int

    Board_blank_count [HEIGHT][WIDTH][2][LEG_INFO_N][4]int
    Notes []string

	callback_begin Func_callback_begin
	callback_count Func_callback_count
	callback_end Func_callback_end
}

type Func_callback_begin func(where Direction)
type Func_callback_count func(k uint, pt Point, where Direction, part Direction) bool
type Func_callback_end func(where Direction) bool


func (self *Bot) Init_data() {
    self.started = false

    self.win_points_num = 0
    self.My_side = WHITE_ID
    self.Your_side = BLACK_ID
    self.Side_this_turn = self.Your_side

    for h:=0; h<HEIGHT; h++ {
        for w:=0; w<WIDTH; w++ {
            self.Board[h][w] = BLANK_ID
            for i:=0; i<LEG_INFO_N; i++ {
                self.Board_leg_type[h][w][int(WHITE_ID)][i] = 1
                self.Board_leg_type[h][w][int(BLACK_ID)][i] = 1
            }
        }
    }
    for h:=0; h<HEIGHT; h++ {
        for w:=0; w<WIDTH; w++ {
            for i:=0; i<LEG_INFO_N; i++ {
                for j:=0; j<4; j++ {
                    self.Board_blank_count[h][w][int(WHITE_ID)][i][j] = 1
                    self.Board_blank_count[h][w][int(BLACK_ID)][i][j] = 1
                }
            }
        }
    }

    self.Notes = self.Notes[:0]
}




func (self *Bot) Notes_dumps() {
    Chess_log(fmt.Sprintf("Notes[%d]: %s", len(self.Notes), strings.ToLower(strings.Join(self.Notes, ""))),
        "INFO")
}


func (self *Bot) Board_string_block() string{
	var board_block_lines []string
	
    board_separate_line := strings.Repeat("- ", WIDTH)
	board_block_lines = append(board_block_lines, "   " + board_separate_line)
	
    for i:=HEIGHT; i>0; i-- {
        var tmp_string []string
        for _, note_info := range self.Board[i-1] {
            tmp_string = append(tmp_string, ID_TO_NOTE[note_info])
        }
		board_block_lines = append(board_block_lines, fmt.Sprintf("%2d|%s|", i, strings.Join(tmp_string, " ")))		
    }
	board_block_lines = append(board_block_lines, "   " + board_separate_line)	
	
    var tmp_labels []string
    for i:=0; i<WIDTH; i++ {
        tmp_labels = append(tmp_labels, idtoa(i))
    }
	board_block_lines = append(board_block_lines, "   " + strings.Join(tmp_labels, " "))

	return strings.Join(board_block_lines, "\n")
}


func (self *Bot) Board_dumps() {
	Chess_log(self.Board_string_block(), "INFO")
}


func (self *Bot) board_debug_dumps() {
    for h:=0; h<HEIGHT; h++ {
        for w:=0; w<WIDTH; w++ {
            if self.get_board_at_point(Point{h, w}) != BLANK_ID {
                continue
            }
            Chess_log(fmt.Sprintf("%s: %s, BLACK: %s", Get_label_of_point(Point{h, w}),
                ID_TO_NOTE[self.Board[h][w]],
                self.Board_leg_type[h][w][BLACK_ID]), "DEBUG")
            Chess_log(fmt.Sprintf("%s: %s, WHITE: %s", Get_label_of_point(Point{h, w}),
                ID_TO_NOTE[self.Board[h][w]],
                self.Board_leg_type[h][w][WHITE_ID]), "DEBUG")
        }
    }
}


func (self *Bot) Board_loads(board_block string) (err error){
    board_block_lines := strings.Split(board_block, "\n")
    if len(board_block_lines) < HEIGHT + 5 {
        return errors.New("board_loads: not enough lines.")
    }

    self.Init_data()
    Chess_log("load board start.", "DEBUG")

    count_balance := 0
    for height, line_side_notes := range board_block_lines[2:len(board_block_lines)-3] {
        line_parts := strings.Split(line_side_notes, "|")
        _, side_notes, _ := line_parts[0], line_parts[1], line_parts[2]
        for i:=0; i<WIDTH; i++ {
            note := string(side_notes[i*2])
            if _, ok := NOTE_TO_ID[note]; ok {
                self.set_board_at_point(Point{HEIGHT-height-1, i}, NOTE_TO_ID[note])
            } else {
                return errors.New(fmt.Sprintf("board_loads: note '%s' is illegal.", note))
            }
            if note == BLACK {
                count_balance += 1
                self.swap_turn_side()
            } else if note == WHITE {
                count_balance -= 1
                self.swap_turn_side()
            }
        }
    }
    if count_balance > 1 || count_balance < 0 {
        return errors.New(fmt.Sprintf("board_loads: notes[%d] is not balance.", count_balance))
    }
    if self.Side_this_turn == self.Your_side {
        // assert(count_balance == 0)
        self.swap_user_side()
    } else {
        //assert(count_balance == 1)
    }
    self.Get_score_of_blanks_for_side(BLACK_ID, true)
    self.Get_score_of_blanks_for_side(WHITE_ID, true)
    Chess_log("load board ok.", "DEBUG")
    return nil
}


func (self *Bot) Light_on_win_points() {
    for i, pt := range self.win_points {
        if i >= self.win_points_num {
            break
        }
        test_side := self.get_board_at_point(pt)
        test_side_win := MARK_WIN[test_side]
        self.set_board_at_point(pt, test_side_win)
    }
}


func (self *Bot) can_put_at_point(pt Point) bool{
    if pt.H < 0 || pt.H >= HEIGHT {
        Chess_log(fmt.Sprintf("point_h(%d) out of range.", pt.H), "DEBUG")
        return false
    }
    if pt.W < 0 || pt.W >= WIDTH {
        Chess_log(fmt.Sprintf("point_w(%d) out of range.", pt.W), "DEBUG")
        return false
    }
    if self.get_board_at_point((pt)) != BLANK_ID {
        Chess_log(fmt.Sprintf("put twice. (%s)", Get_label_of_point(pt)), "DEBUG")
        return false
    }
    return true
}


func (self *Bot) swap_user_side() {
    self.My_side, self.Your_side = self.Your_side, self.My_side
}


func (self *Bot) swap_turn_side() {
    if self.Side_this_turn == self.My_side {
        self.Side_this_turn = self.Your_side
    } else {
        self.Side_this_turn = self.My_side
    }
}


func (self *Bot) get_board_at_point(pt Point) GoSide{
    return self.Board[pt.H][pt.W]
}


func (self *Bot) set_board_at_point(pt Point, side_note GoSide) {
    self.Board[pt.H][pt.W] = side_note
}


func (self *Bot) Put_chessman_at_point(put_side GoSide, pt Point) (err error) {
    if self.Side_this_turn != put_side {
        return errors.New(fmt.Sprintf("not %s turn.", ID_TO_NOTE[put_side]))
    }
    if self.can_put_at_point(pt) {
        self.set_board_at_point(pt, self.Side_this_turn)
        self.update_put_around_point(pt, self.Side_this_turn)
        operate := fmt.Sprintf("%s %s %s", OP_PUT, Get_label_of_point(pt), ID_TO_NOTE[self.Side_this_turn])
        Chess_operate(operate)
        self.Notes = append(self.Notes, Get_label_of_point(pt))
        self.swap_turn_side()
        return nil
    }
    return nil
}


func (self *Bot) Get_point_of_chessman(get_side GoSide) (pt Point, err error) {
    if self.Side_this_turn != get_side {
        Chess_log(fmt.Sprintf("not %s turn.", ID_TO_NOTE[get_side]), "DEBUG")
        return pt, errors.New("get_point_of_chessman fail.")
    }

	var line string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = strings.ToUpper(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
	Chess_log(fmt.Sprintf("get: %v", line), "INFO")
	
    if line == "START" {
        if !self.started {
            self.started = true
            self.swap_user_side()
            self.Side_this_turn = self.My_side
        }
        return pt, errors.New("get_point_of_chessman double start.")
    }
    if !strings.HasPrefix(line, OP_PUT) {
        return pt, errors.New("get_point_of_chessman not put.")
    }

    tmp_line_split := make([]string, 3)
    if len(strings.Split(line, " ")) == 2 {
        tmp_line_split = strings.Split(line, " ")
    } else {
        tmp_line_split = strings.SplitN(line, " ", 2)
    }
	Chess_log(fmt.Sprintf("here: %v", tmp_line_split), "INFO")
    op_token, point_token := tmp_line_split[0], tmp_line_split[1]
    point_w_s, point_h_s := point_token[0], point_token[1:]
    if point_h, err := strconv.Atoi(point_h_s); err == nil {
        point_w := atoid(string(point_w_s))
        pt = Point{point_h-1, point_w}
    } else {
        return pt, errors.New("get_point_of_chessman can't parse.")
    }
    if op_token == OP_PUT && self.can_put_at_point(pt) {
        self.started = true
        self.set_board_at_point(pt, self.Side_this_turn)
        self.update_put_around_point(pt, self.Side_this_turn)
        self.Notes = append(self.Notes, Get_label_of_point(pt))
        self.swap_turn_side()
        return pt, nil
    }
    return pt, errors.New("get_point_of_chessman can't put.")
}


func min2(a int, b int) int{
    if a < b {
        return a
    }
    return b
}

func min3(a int, b int, c int) int{
    return min2(min2(a, b), c)
}



func (self *Bot) detect_positions_around_point(test_pt Point, test_side GoSide) bool{
    // 包括回调函数 {
    //
    // self.callback_begin
    // self.callback_count
    // self.callback_end
    //
    // 在遍历中心点的4个主要方向时回调.
    // 开始遍历此方向时调用一次: callback_begin,
    // 遍历此方向上每一个位置调用: callback_count
    // 此方向遍历完毕调用: callback_end
    // 如果callback_end返回true, 则函数提前返回true
    //
    // 遍历完8个方向后, 返回false
    //

    self.center_point = test_pt
    self.center_side = self.get_board_at_point(test_pt)
    self.test_side = test_side

    h, w := test_pt.H, test_pt.W
    // test row (-)
    self.callback_begin(DIRECTION_ROW)
    for k:=0; k<min2(WIN_NUM-1, w); k++ {
        pt := Point{h, w-k-1}
        if self.callback_count(uint(k), pt, DIRECTION_ROW, DIRECTION_WEST) {
            break
        }
    }
    for k:=0; k<min2(WIN_NUM-1, WIDTH-w-1); k++ {
        pt := Point{h, w+k+1}
        if self.callback_count(uint(k), pt, DIRECTION_ROW, DIRECTION_EAST) {
            break
        }
    }
    if self.callback_end(DIRECTION_ROW) {
        return true
    }
    // test col (|)
    self.callback_begin(DIRECTION_COL)
    for k:=0; k<min2(WIN_NUM-1, h); k++ {
        pt := Point{h-k-1, w}
        if self.callback_count(uint(k), pt, DIRECTION_COL, DIRECTION_SOUTH) {
            break
        }
    }
    for k:=0; k<min2(WIN_NUM-1, WIDTH-h-1); k++ {
        pt := Point{h+k+1, w}
        if self.callback_count(uint(k), pt, DIRECTION_COL, DIRECTION_NORTH) {
            break
        }
    }
    if self.callback_end(DIRECTION_COL) {
        return true
    }
    // test down (\)
    self.callback_begin(DIRECTION_DOWN)
    for k:=0; k<min3(WIN_NUM-1, HEIGHT-h-1, w); k++ {
        pt := Point{h+k+1, w-k-1}
        if self.callback_count(uint(k), pt, DIRECTION_DOWN, DIRECTION_NORTHWEST) {
            break
        }
    }
    for k:=0; k<min3(WIN_NUM-1, h, WIDTH-w-1); k++ {
        pt := Point{h-k-1, w+k+1}
        if self.callback_count(uint(k), pt, DIRECTION_DOWN, DIRECTION_SOUTHEAST) {
            break
        }
    }
    if self.callback_end(DIRECTION_DOWN) {
        return true
    }
    // test up (/)
    self.callback_begin(DIRECTION_UP)
    for k:=0; k<min3(WIN_NUM-1, h, w); k++ {
        pt := Point{h-k-1, w-k-1}
        if self.callback_count(uint(k), pt, DIRECTION_UP, DIRECTION_SOUTHWEST) {
            break
        }
    }
    for k:=0; k<min3(WIN_NUM-1, HEIGHT-h-1, WIDTH-w-1); k++ {
        pt := Point{h+k+1, w+k+1}
        if self.callback_count(uint(k), pt, DIRECTION_UP, DIRECTION_NORTHEAST) {
            break
        }
    }
    if self.callback_end(DIRECTION_UP) {
        return true
    }
    return false
}


func (self *Bot) callback_begin_winner(where Direction) {
	self.win_points_num = 0	
    self.win_points[self.win_points_num] = self.center_point
    self.win_points_num++
}
func (self *Bot) callback_count_winner(k uint, pt Point, where Direction, part Direction) bool{
    if self.get_board_at_point(pt) != self.test_side {
        return true
    }
    self.win_points[self.win_points_num] = pt
    self.win_points_num++
    return false
}
func (self *Bot) callback_end_winner(where Direction) bool{
    if self.win_points_num >= WIN_NUM {
        return true
    }
    return false
}
func (self *Bot) Is_winner(test_side GoSide) bool{
    self.callback_count = self.callback_count_winner
    self.callback_end = self.callback_end_winner
    self.callback_begin = self.callback_begin_winner

    for h:=0; h<HEIGHT; h++ {
        for w:=0; w<WIDTH; w++ {
            if self.get_board_at_point(Point{h, w}) != test_side {
                continue
            }
            if self.detect_positions_around_point(Point{h, w}, test_side) {
                return true
            }
        }
    }
    return false
}


func (self *Bot) Win_test(pt Point, test_side GoSide) bool{
    self.callback_count = self.callback_count_winner
    self.callback_end = self.callback_end_winner
    self.callback_begin = self.callback_begin_winner

    self.set_board_at_point(pt, test_side)
    if self.detect_positions_around_point(pt, test_side) {
        self.set_board_at_point(pt, BLANK_ID)
        return true
    }
    self.set_board_at_point(pt, BLANK_ID)
    return false
}


func (self *Bot) callback_begin_legtype(where Direction) {
    // 空点汇合计数
    return
}
func (self *Bot) callback_count_legtype(k uint, pt Point, where Direction, part Direction) bool{
    k_type := 1<<k
    legtype := self.direction_legtype_count[part]
    legtype += k_type

    if self.get_board_at_point(pt) == OPPOSITE_SIDE[self.test_side] {
        return true
    }
    if self.get_board_at_point(pt) != BLANK_ID {
        legtype += k_type
    }
    self.direction_legtype_count[part] = legtype
    return false
}
func (self *Bot) callback_end_legtype(where Direction) bool{
    return false
}


func (self *Bot) callback_begin_update_blank_score_after_put(where Direction) {
    // 放子更新
    return
}
func (self *Bot) callback_count_update_blank_score_after_put(k uint, pt Point, where Direction, part Direction) bool{
    if self.get_board_at_point(pt) != BLANK_ID {
        return false
    }
    opposite_part := OPPOSITE_DIRECTION[part]

    self.Board_blank_count[pt.H][pt.W][int(BLACK_ID)][opposite_part][k] = self.Board_leg_type[pt.H][pt.W][int(BLACK_ID)][opposite_part]
    self.Board_blank_count[pt.H][pt.W][int(WHITE_ID)][opposite_part][k] = self.Board_leg_type[pt.H][pt.W][int(WHITE_ID)][opposite_part]

    k_type := 1<<k
    temp_c := self.Board_leg_type[pt.H][pt.W][self.test_side][opposite_part]
    if temp_c > k_type {
        temp_c += k_type
        self.Board_leg_type[pt.H][pt.W][self.test_side][opposite_part] = temp_c
    }
    temp_c = self.Board_leg_type[pt.H][pt.W][OPPOSITE_SIDE[self.test_side]][opposite_part]
    if temp_c > k_type {
        temp_c &= (k_type-1)
        temp_c += k_type
        self.Board_leg_type[pt.H][pt.W][OPPOSITE_SIDE[self.test_side]][opposite_part] = temp_c
    }
    self.update_total_score(pt, BLACK_ID)
    self.update_total_score(pt, WHITE_ID)
    return false
}
func (self *Bot) callback_end_update_blank_score_after_put(where Direction) bool{
    return false
}


func (self *Bot) update_total_score(pt Point, test_side GoSide) {
    total_score := 0
    total_dup_score := 0
    LEG_TYPE_TO_COUNT := []int{
        0, 0, 0, 1,
        0, 1, 1, 2,
        0, 1, 1, 2,
        1, 2, 2, 3,
        0, 1, 1, 2,
        1, 2, 2, 3,
        1, 2, 2, 3,
        2, 2, 3, 4,}
    LEG_TYPE_TO_SPACE := []int{
        0, 0, 1, 1,
        2, 2, 2, 2,
        3, 3, 3, 3,
        3, 3, 3, 3,
        4, 4, 4, 4,
        4, 4, 4, 4,
        4, 4, 4, 4,
        4, 4, 4, 4,}
    LEG_TYPE_TO_ADDITION_COUNT := []int{
        0, 0, 0, 0,
        0, 0, 0, 2,
        0, 0, 0, 2,
        0, 0, 0, 3,
        0, 0, 0, 2,
        0, 0, 0, 3,
        0, 0, 0, 2,
        0, 0, 0, 4,}
    LEG_TYPE_TO_DOUBLE_ADDITION_COUNT := []int{
        0, 0, 0, 1,
        0, 1, 0, 2,
        0, 1, 0, 2,
        0, 1, 0, 3,
        0, 1, 0, 2,
        0, 1, 0, 3,
        0, 1, 0, 2,
        0, 1, 0, 4,}

    dir_types := self.Board_leg_type[pt.H][pt.W][test_side][:LEG_INFO_N-3]
    for direction, t := range dir_types {
        opposite_dir := OPPOSITE_DIRECTION[Direction(direction)]
        opposite_t := dir_types[opposite_dir]
        // test space
        t_score := LEG_TYPE_TO_SPACE[t]
        opposite_t_score := LEG_TYPE_TO_SPACE[opposite_t]
        if t_score + opposite_t_score < WIN_NUM-1 {
            continue
        }
        // base score
        total_score += LEG_TYPE_TO_COUNT[t]
        total_dup_score += LEG_TYPE_TO_ADDITION_COUNT[t]

        // addition score
        t_score = LEG_TYPE_TO_DOUBLE_ADDITION_COUNT[t]
        opposite_t_score = LEG_TYPE_TO_DOUBLE_ADDITION_COUNT[opposite_t]
        if t_score > 0 && opposite_t_score > 0 {
            total_dup_score += t_score
        }
    }
    self.Board_leg_type[pt.H][pt.W][test_side][LEG_INFO_IDX_DUP_SUM_SCORE] = total_score + total_dup_score
    self.Board_leg_type[pt.H][pt.W][test_side][LEG_INFO_IDX_SUM_SCORE] = total_score
}


func (self *Bot) callback_begin_update_blank_score_after_remove(where Direction) {
    // 放子更新
    return
}
func (self *Bot) callback_count_update_blank_score_after_remove(k uint, pt Point, where Direction, part Direction) bool{
    if self.get_board_at_point(pt) != BLANK_ID {
        return false
    }
    opposite_part := OPPOSITE_DIRECTION[part]
    self.Board_leg_type[pt.H][pt.W][int(BLACK_ID)][opposite_part] = self.Board_blank_count[pt.H][pt.W][int(BLACK_ID)][opposite_part][k]
    self.Board_leg_type[pt.H][pt.W][int(WHITE_ID)][opposite_part] = self.Board_blank_count[pt.H][pt.W][int(WHITE_ID)][opposite_part][k]
    self.update_total_score(pt, BLACK_ID)
    self.update_total_score(pt, WHITE_ID)
    return false
}
func (self *Bot) callback_end_update_blank_score_after_remove(where Direction) bool{
    return false
}


// func (self *Bot) callback_begin(where Direction) {
//     return
// }
// func (self *Bot) callback_count(k uint, pt Point, where Direction, part Direction) bool{
//     return false
// }
// func (self *Bot) callback_end(where Direction) bool{
//     return false
// }


func (self *Bot) update_remove_around_point(pt Point) {
    self.callback_count = self.callback_count_update_blank_score_after_remove
    self.callback_end = self.callback_end_update_blank_score_after_remove
    self.callback_begin = self.callback_begin_update_blank_score_after_remove
    self.detect_positions_around_point(pt, BLANK_ID)
}


func (self *Bot) update_put_around_point(pt Point, test_side GoSide) {
    self.callback_count = self.callback_count_update_blank_score_after_put
    self.callback_end = self.callback_end_update_blank_score_after_put
    self.callback_begin = self.callback_begin_update_blank_score_after_put
    self.detect_positions_around_point(pt, test_side)
}


func (self *Bot) Get_score_of_blanks_for_side(test_side GoSide, is_dup_enforce bool) PairList {
    // test_side的每一个棋子, 对它米子型中心WIN_NUM范围的空白的位置贡献记分
    // 返回所有的空白位置坐标和对应的累计记分, pair

    // 获取所有在棋盘中test_side棋子的位置坐标
    all_my_blank_points_count := map[Point]int{}
    for h:=0; h<HEIGHT; h++ {
        for w:=0; w<WIDTH; w++ {
            pt := Point{h, w}
            if self.get_board_at_point(pt) != BLANK_ID {
                continue
            }
            self.direction_legtype_count = self.Board_leg_type[pt.H][pt.W][test_side][:]

            if self.direction_legtype_count[POINT_NEED_UPDATE] == 0 {
				if self.direction_legtype_count[LEG_INFO_IDX_SUM_SCORE] > 0	{
					Chess_log(fmt.Sprintf("%s GET SCORE[%s]: %v", ID_TO_NOTE[test_side], Get_label_of_point(pt),
						self.direction_legtype_count), "DEBUG")
				}
			} else {
                for i:=0; i<LEG_INFO_N; i++ {
                    self.direction_legtype_count[i] = 1
                }
                self.direction_legtype_count[POINT_NEED_UPDATE] = 0

                self.callback_count = self.callback_count_legtype
                self.callback_end = self.callback_end_legtype
                self.callback_begin = self.callback_begin_legtype
                self.detect_positions_around_point(pt, test_side)

                self.update_total_score(pt, test_side)
                Chess_log(fmt.Sprintf("%s UPDATE SCORE[%s]: %v", ID_TO_NOTE[test_side], Get_label_of_point(pt),
                    self.direction_legtype_count), "DEBUG")
            }
			blank_score := 0
            if is_dup_enforce {
                blank_score = self.direction_legtype_count[LEG_INFO_IDX_DUP_SUM_SCORE]
            } else {
                blank_score = self.direction_legtype_count[LEG_INFO_IDX_SUM_SCORE]
            }
            if blank_score > 0 {
                all_my_blank_points_count[pt] = blank_score
            }
        }
    }

	all_my_blank_points_count_pair := Rank_by_point_count(all_my_blank_points_count)
	if G_debug_info {
		tmp_labels := []string{}
		for _, ppair := range all_my_blank_points_count_pair {
			if ppair.Value > 0 {
				tmp_labels = append(tmp_labels, fmt.Sprintf("%s:%d", Get_label_of_point(ppair.Key), ppair.Value))
			}
		}
		Chess_log(fmt.Sprintf("%s Score: %s", ID_TO_NOTE[test_side], strings.Join(tmp_labels, ", ")), "DEBUG")
	}
    return all_my_blank_points_count_pair
}


func (self *Bot) Is_a_good_choice(choice_pt Point, my_side GoSide, your_side GoSide, max_level int) bool{
    // todo: 层序遍历, 最高得分先检查
    if max_level == 0 {
        return false
    }

    self.set_board_at_point(choice_pt, my_side)
    self.update_put_around_point(choice_pt, my_side)
    Chess_log(fmt.Sprintf("%s TEST GOOD CHOICE[%d]: %s", ID_TO_NOTE[my_side], max_level,
        Get_label_of_point(choice_pt)), "DEBUG")

    is_dup_enforce := true
    all_my_blank_points_count_pair := self.Get_score_of_blanks_for_side(my_side, is_dup_enforce)
    count_win_point := 0
    for _, ppair := range all_my_blank_points_count_pair {
		my_pt, count := ppair.Key, ppair.Value		
        // 先扫一遍有没有多处直接胜利的, count<4的点不可能胜利
        if count < 4 {
            continue
        }
        if self.Win_test(my_pt, my_side) {
            count_win_point += 1
            if count_win_point > 1 {
                self.update_remove_around_point(choice_pt)
                self.set_board_at_point(choice_pt, BLANK_ID)
                return true
            }
        }
    }
    tested_not_good_pt := map[Point]bool{}
    for _, ppair := range all_my_blank_points_count_pair {
		my_pt, count := ppair.Key, ppair.Value		
        if count < 1 {
            continue
        }
		tested_not_good_pt[my_pt] = true
        if is_bad := self.Is_a_bad_choice(my_pt, your_side, my_side, max_level); !is_bad {
            self.update_remove_around_point(choice_pt)
            self.set_board_at_point(choice_pt, BLANK_ID)
            return false
        }
    }
    is_dup_enforce = true
    all_your_blank_points_count_pair := self.Get_score_of_blanks_for_side(your_side, is_dup_enforce)
    for _, ppair := range all_your_blank_points_count_pair {
		your_pt, count := ppair.Key, ppair.Value		
        if count < 1 {
            continue
        }
        if tested_not_good_pt[your_pt] {
            continue
        }
        if is_bad := self.Is_a_bad_choice(your_pt, your_side, my_side, max_level); !is_bad {
            Chess_log(fmt.Sprintf("%s GET BAD CHOICE[%d]: %s",
                ID_TO_NOTE[your_side], max_level-1,
                Get_label_of_point(your_pt)), "DEBUG")

            self.update_remove_around_point(choice_pt)
            self.set_board_at_point(choice_pt, BLANK_ID)
            return false
        }
    }
    self.update_remove_around_point(choice_pt)
    self.set_board_at_point(choice_pt, BLANK_ID)
    return true
}


func (self *Bot) Is_a_bad_choice(choice_pt Point, my_side GoSide, your_side GoSide, max_level int) bool{
    // todo: 层序遍历, 最高得分先检查
    if max_level == 0 {
        return false
    }

    self.set_board_at_point(choice_pt, my_side)
    self.update_put_around_point(choice_pt, my_side)
    Chess_log(fmt.Sprintf("%s TEST BAD CHOICE[%d]: %s", ID_TO_NOTE[my_side], max_level,
        Get_label_of_point(choice_pt)), "DEBUG")

    is_dup_enforce := true
    all_your_blank_points_count_pair := self.Get_score_of_blanks_for_side(your_side, is_dup_enforce)
    for _, ppair := range all_your_blank_points_count_pair {
		your_pt, count := ppair.Key, ppair.Value		
        // 先扫一遍有没有直接胜利的, count<4的点不可能胜利
        if count >= 4 {
            if self.Win_test(your_pt, your_side) {
                self.update_remove_around_point(choice_pt)
                self.set_board_at_point(choice_pt, BLANK_ID)
                return true
            }
        }
    }

    for _, ppair := range all_your_blank_points_count_pair {
		your_pt, count := ppair.Key, ppair.Value		
        if count > 2 {
            // tofix: 不应该忽视count==1的点, 但为了减少计算
            if self.Is_a_good_choice(your_pt, your_side, my_side, max_level-1) {
                Chess_log(fmt.Sprintf("%s GET GOOD CHOICE[%d]: %s",
                    ID_TO_NOTE[your_side], max_level-1,
                    Get_label_of_point(your_pt)), "DEBUG")

                self.update_remove_around_point(choice_pt)
                self.set_board_at_point(choice_pt, BLANK_ID)
                return true
            }
        }
    }

    is_dup_enforce = true
    all_my_blank_points_count_pair := self.Get_score_of_blanks_for_side(my_side, is_dup_enforce)
    for _, ppair := range all_my_blank_points_count_pair {
		my_pt, count := ppair.Key, ppair.Value		
        if count > 2 {
            // tofix: 不应该忽视count==1的点, 但为了减少计算
            if self.Is_a_good_choice(my_pt, your_side, my_side, max_level-1) {
                Chess_log(fmt.Sprintf("%s GET GOOD CHOICE[%d]: %s",
                    ID_TO_NOTE[your_side], max_level-1,
                    Get_label_of_point(my_pt)), "DEBUG")

                self.update_remove_around_point(choice_pt)
                self.set_board_at_point(choice_pt, BLANK_ID)
                return true
            }
        }
    }

    self.update_remove_around_point(choice_pt)
    self.set_board_at_point(choice_pt, BLANK_ID)
    return false
}
