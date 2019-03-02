package main

import "testing"

func TestVerifyAnswers(t *testing.T) {
	quiz := Quiz{
		QuizItem{"5+5", "10"},
		QuizItem{"7+3", "10"},
		QuizItem{"1+1", "2"},
	}
	result := verifyAnswers([]string{"10", "10", "2"}, quiz)
	if result.Score != 3 {
		t.Errorf("verifyAnswers() = %d; want 1", result.Score)
	}
}
