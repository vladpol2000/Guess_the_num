package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var d int // Переменная для числа, вводимого игроком

// Структура для сохранения статистики игр в json файл
type GameResult struct {
	Date     string `json:"date"`
	Result   string `json:"result"`
	Attempts int    `json:"attempts"`
}

func main() {
	var (
		rand_num    int
		level       int
		play_again  bool = true
		num_of_game int  = 1
		answer      string
		max_level   int    = 3
		tips_text   string = "Введите номер уровня игры. Easy - 1, Medium - 2 или Hard - 3: "
	)

	// Используем мапу для уровней игры. Значением по ключу является массив с целыми числами: 1-ое число - диапазон, 2 число - кол-во попыток.
	levels := map[int][]int{
		1: {50, 15},
		2: {100, 10},
		3: {200, 5},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for play_again {
		if num_of_game == 1 {
			fmt.Print("Добро пожаловать в игру " + color.HiGreenString("\"Угадай число\"\n"))
			fmt.Print(tips_text)

			level = validateInput(scanner, max_level, tips_text)

			// Генерация случайного числа от 0 до 100
			rand_num = rand.Intn(levels[level][0] + 1)

			fmt.Printf("Давайте сыграем в игру. Попробуйте отгадать загаданное число от 1 до %d. У вас будет всего лишь %d попыток. \n", levels[level][0], levels[level][1])

			msg, result := play_game(levels[level][1], rand_num, levels[level][0])
			fmt.Println(msg)
			saveResult(result)

			num_of_game++
		} else {
			if repeatGame(answer) {
				fmt.Print(tips_text)

				level := validateInput(scanner, max_level, tips_text)

				// Генерация случайного числа от 0 до 100
				rand_num = rand.Intn(levels[level][0] + 1)

				fmt.Printf("Давайте сыграем в игру. Попробуйте отгадать загаданное число от 1 до %d. У вас будет всего лишь %d попыток. \n", levels[level][0], levels[level][1])

				msg, result := play_game(levels[level][1], rand_num, levels[level][0])
				fmt.Println(msg)
				saveResult(result)

				num_of_game++
			} else {
				fmt.Println("Спасибо за игру!")
				return
			}
		}
	}
}

func play_game(nums_of_tries int, rand_num int, max_num int) (string, GameResult) {
	list_nums := make([]int, 0) // Переменная для введённых игроком чисел
	scanner := bufio.NewScanner(os.Stdin)
	var i int = 0
	tips_text := "Введите загаданное число: "

	// Цикл для n попыток сыграть в игру
	for i < nums_of_tries {

		fmt.Print(tips_text)

		// Выполняем проверку, что введено число
		d = validateInput(scanner, max_num, tips_text)

		list_nums = append(list_nums, d)

		fmt.Printf("Вы уже вводили следующие числа %v \n", list_nums)

		if d == rand_num {
			return color.GreenString("Вы угадали!🙌"), GameResult{
				Date:     time.Now().Format(time.RFC3339),
				Result:   "Победа",
				Attempts: i + 1,
			}
		} else {
			switch {
			case d > rand_num:
				fmt.Println("Секретное число меньше👇" + tips(rand_num, d))
			case d < rand_num:
				fmt.Println("Секретное число больше👆" + tips(rand_num, d))
			case d == rand_num:
			}
		}
		i++
	}
	res := color.RedString("Вы проиграли!😢\n") + "Секретное число было: " + strconv.Itoa(rand_num)
	return res, GameResult{
		Date:     time.Now().Format(time.RFC3339),
		Result:   "Проигрыш",
		Attempts: nums_of_tries,
	}
}

func tips(rand_num int, d int) string {
	switch {
	case absInt(d-rand_num) <= 5:
		return (" Ты совсем близок 🔥")
	case absInt(d-rand_num) <= 15 && absInt(d-rand_num) > 5:
		return (" Уже теплее 🙂")
	default:
		return (" Тут очень холодно ❄️")
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func saveResult(result GameResult) {
	var results []GameResult

	// Читаем файл, если он есть
	file, err := os.ReadFile("result.json")
	if err == nil {
		json.Unmarshal(file, &results)
	}

	// Добавляем новый результат
	results = append(results, result)

	// Перезаписываем файл
	data, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("result.json", data, 0644)
}

func validateInput(scanner *bufio.Scanner, max_num int, tips_text string) int {
	for {
		scanner.Scan()
		input := scanner.Text()

		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка ввода. В консоль можно вводить только числа")
			fmt.Print(tips_text)
			continue
		} else if num > max_num {
			fmt.Println("Ошибка ввода. Введённое число не может превышать максимальное -", max_num)
			fmt.Print(tips_text)
			continue
		}

		return num
	}
}

func repeatGame(answer string) bool {
	for {
		fmt.Println("Хотите сыграть ещё одну игру? Введите Да или Нет: ")
		fmt.Scan(&answer)

		answer = strings.ToLower(answer)

		if answer == "да" {
			return true
		} else if answer == "нет" {
			return false
		} else {
			fmt.Println("Вы ввели неверное значение. Попробуйте ещё раз")
			continue
		}
	}
}
