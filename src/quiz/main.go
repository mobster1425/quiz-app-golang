package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// number of questions to ask the user
// const totalQuestions = 5

type Question struct {
	question string
	answer   string
}

// QuestionSet represents a set of quiz questions.
/*
type QuestionSet struct {
	questions []Question
}
*/

func main() {

	var qs Question
	//var q QuestionSet
	filename, timeLimit, shuffle := qs.readArguements()
	f, err := qs.openFile(filename)
	if err != nil {
		return
	}
	questions, err := qs.readCSV(f)

	if err != nil {
		// err := fmt.Errorf("Error in Reading Questions")
		fmt.Println(err.Error())
		return
	}

	if questions == nil {
		return
	}
	if shuffle {
		shuffleQuestions(questions)
	}

	score, err := askquestion(questions, timeLimit)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Your Score %d/%d\n", score, len(questions))
}

func (q *Question) readArguements() (string, int, bool) {
	fileName := flag.String("filename", "quiz-problem.csv", "the csv that contains the quizes")
	//default value is 30 and default value for filename is quiz-problem.csv
	timeLimit := flag.Int("limit", 30, "time limit for each question")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz order")
	flag.Parse()
	return *fileName, *timeLimit, *shuffle
}

// readCSV is a method of the Question struct.
func (q *Question) readCSV(f io.Reader) ([]Question, error) {

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t' // Set the delimiter to tab character
	allQuestions, err := csvReader.ReadAll()

	// allQuestions, err := csv.NewReader(f).ReadAll()

	if err != nil {
		return nil, err
	}

	numOfQues := len(allQuestions)
	if numOfQues == 0 {
		return nil, fmt.Errorf("No Question in file")
	}

	var data []Question
	for _, line := range allQuestions {
		if len(line) < 2 {
			return nil, fmt.Errorf("Invalid format for question and answer")
		}
		ques := Question{}
		ques.question = line[0]
		ques.answer = line[1]
		data = append(data, ques)
	}

	return data, nil
}

func (q Question) openFile(fileName string) (io.Reader, error) {
	return os.Open(fileName)
}

// get input function is not a method
func getInput(input chan string) {
	for {
		in := bufio.NewReader(os.Stdin)
		result, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input <- result
	}
}

// AskQuestion conducts the quiz and returns the total score and any error encountered.
// it is a method
func askquestion(questions []Question, timeLimit int) (int, error) {
	totalScore := 0
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	done := make(chan string)

	/*
		after starting the goroutine with go getInput(done), the method continues its execution and runs the quiz loop. The goroutine
		running getInput will listen for input from the user in the background and send the input to the done channel when it receives it.

		By using concurrency in this way, the askquestion method can listen for user input while simultaneously
		conducting the quiz with the timer. If the timer expires or the quiz ends, the done channel can be closed to
		 terminate the getInput goroutine,
		ensuring it doesn't run indefinitely. This is how the input and the quiz process can happen concurrently.

	*/
	go getInput(done)

	for i := range questions {
		ans, err := eachQuestion(questions[i].question, questions[i].answer, timer.C, done)
		if err != nil && ans == -1 {
			return totalScore, nil
		}
		totalScore += ans
	}

	return totalScore, nil
}

// not a method

func eachQuestion(Quest string, answer string, timer <-chan time.Time, done <-chan string) (int, error) {
	fmt.Printf("%s: ", Quest)

	for {
		select {
		case <-timer:
			return -1, fmt.Errorf("Time out")
		case ans := <-done:
			score := 0
			trimmedAns := strings.TrimSpace(strings.ToLower(ans))
			trimmedCorrectAns := strings.TrimSpace(strings.ToLower(answer))

			if trimmedAns == trimmedCorrectAns {
				score = 1
			} else {
				return 0, fmt.Errorf("Wrong Answer")
			}

			return score, nil
		}
	}
}

func shuffleQuestions(questions []Question) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
}
