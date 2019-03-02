//TODO implement rest of the option exercises
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func shuffleQuiz(quiz Quiz) {
	rand.Shuffle(len(quiz), func(i, j int) {
		quiz[i], quiz[j] = quiz[j], quiz[i]
	})
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func checkFile(e error, csvFilename *string) {
	if e != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
}

func checkCsvParse(err error) {
	if err != nil {
		exit("Failed to parse the CSV file")
	}
}

/*
QuizItem stores one question and the respective answer
*/
type QuizItem struct {
	Question      string
	CorrectAnswer string
}

/*
Quiz stores the entire set of questions (QuizItems)
*/
type Quiz []QuizItem

/*
QuizItemResult stores the provided answer and wheter it is correct for quiz item
*/
type QuizItemResult struct {
	quizItem       QuizItem
	ProvidedAnswer string
	Correct        bool
}

/*
QuizResult stores the result for the entire quiz and the corresponding Score
 (nr of questions with correct answer)
*/
type QuizResult struct {
	quizResult []QuizItemResult
	Score      int
}

func createQuiz(records [][]string) Quiz {
	//	var quiz Quiz
	quiz := make(Quiz, len(records))
	for i, r := range records {
		quiz[i] = QuizItem{r[0], strings.TrimSpace(r[1])}
	}
	return quiz
}

func verifyAnswers(answers []string, quiz Quiz) QuizResult {
	var quizItemResult []QuizItemResult
	score := 0
	for i, a := range answers {
		isAnswerCorrect := quiz[i].CorrectAnswer == a
		if isAnswerCorrect {
			score = score + 1
		}
		quizItemResult = append(quizItemResult, QuizItemResult{quiz[i], a, isAnswerCorrect})
	}
	for i := len(answers); i < len(quiz); i++ {
		quizItemResult = append(quizItemResult, QuizItemResult{quiz[i], "", false})
	}
	return QuizResult{quizItemResult, score}
}

func askQuestions(quiz Quiz, ch chan string, endChannel chan int) {
	answerScanner := bufio.NewScanner(os.Stdin)
	for _, q := range quiz {
		fmt.Println("Question:", q.Question)
		answerScanner.Scan()
		ch <- answerScanner.Text()
	}
	close(ch)
	endChannel <- 1
}

func receiveAnswers(ch chan string) []string {
	answers := make([]string, len(ch))
	i := 0
	for answer := range ch {
		answers[i] = strings.TrimSpace(answer)
		i++
	}
	return answers
}

func main() {
	// read csv file

	var fileName = flag.String("file", "quiz.csv", "path to csv file containing question, answer records")
	var timeout = flag.Int("limit", 30, "duration of the test in secods")
	flag.Parse()

	f, err := os.Open(*fileName)
	checkFile(err, fileName)
	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	checkCsvParse(err)

	//store in data structure
	quiz := createQuiz(records)
	shuffleQuiz(quiz)
	// ask user for answers
	fmt.Println("You have", timeout, "seconds to complete the quiz")
	fmt.Println("Press any key to start quiz...")
	fmt.Scanln()
	ch := make(chan string, len(quiz))
	endChannel := make(chan int)
	go askQuestions(quiz, ch, endChannel)
	select {
	case <-endChannel:
	case <-time.After(time.Duration(*timeout) * time.Second):
		fmt.Println("Quiz time limit reached!")
		close(ch)
	}
	answers := receiveAnswers(ch)
	// check answers and apply Score
	result := verifyAnswers(answers, quiz)
	fmt.Println("Correct:", result.Score, "questions out of ", len(quiz))
	//fmt.Println(result)
}
