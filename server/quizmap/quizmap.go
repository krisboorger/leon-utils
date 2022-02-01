package quizmap

import (
	"errors"
	"reflect"
)

type void struct{}

var setMember void

type StringSet map[string]void
type AnswerMap map[string]StringSet
type QuestionMap map[string]*Question

type Question struct {
	SubQuestions QuestionMap
	Answers      AnswerMap
}

func deepElem(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	return value
}

var ErrAnswerMalformed = errors.New("answer map malformed")
var ErrNotFinalSubquestion = errors.New("not final subquestion - can't edit answers")
var ErrFinalSubquestion = errors.New("final subquestion - can't edit it's subquestions")

func (q *Question) Update(attempt string, newAnswers interface{}) error {
	value := deepElem(reflect.ValueOf(newAnswers))
	switch value.Kind() {
	case reflect.Map:
		if q.Answers != nil {
			return ErrFinalSubquestion
		}
		if q.SubQuestions == nil {
			q.SubQuestions = QuestionMap{}
		}
		for _, subQNameValue := range value.MapKeys() {
			if subQNameValue.Kind() != reflect.String {
				return ErrAnswerMalformed
			}

			subQName := subQNameValue.String()
			subQAnswers := value.MapIndex(subQNameValue).Interface()

			subQuestion := q.SubQuestions[subQName]
			if subQuestion == nil {
				subQuestion = &Question{}
				q.SubQuestions[subQName] = subQuestion
			}
			err := subQuestion.Update(attempt, subQAnswers)
			if err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		if q.SubQuestions != nil {
			return ErrNotFinalSubquestion
		}
		if q.Answers == nil {
			q.Answers = AnswerMap{}
		}

		// opt out from previous answers
		for _, v := range q.Answers {
			delete(v, attempt)
		}

		// opt in for the selected answers
		for i := 0; i < value.Len(); i++ {
			answerValue := deepElem(value.Index(i))
			if answerValue.Kind() != reflect.String {
				return ErrAnswerMalformed
			}
			answer := answerValue.String()

			if answer == "" {
				continue
			}
			stringSet := q.Answers[answer]
			if stringSet == nil {
				stringSet = make(StringSet)
				q.Answers[answer] = stringSet
			}
			stringSet[attempt] = setMember
		}
	default:
		return ErrAnswerMalformed
	}
	return nil
}

func (q *Question) ToCounts() map[string]interface{} {
	result := make(map[string]interface{})
	if q.SubQuestions != nil {
		// has sub question(s), call ToCounts on them
		for subQName, subQ := range q.SubQuestions {
			result[subQName] = subQ.ToCounts()
		}
		return result
	} else if q.Answers != nil {
		// last sub question, return answers
		for answer, answerSet := range q.Answers {
			result[answer] = len(answerSet)
		}
		return result
	}
	// both are nil - dangling question
	return result
}

type QuizMap map[string]QuestionMap

func (q QuizMap) UpdateAnswer(quiz string, attempt string, questionName string, newAnswers interface{}) error {
	questionMap := q[quiz]
	if questionMap == nil {
		questionMap = make(QuestionMap)
		q[quiz] = questionMap
	}
	question := questionMap[questionName]
	if question == nil {
		question = &Question{}
		questionMap[questionName] = question
	}

	return question.Update(attempt, newAnswers)
}

func (q QuizMap) GetAnswerCounts(quiz string, questionNames []string) map[string]interface{} {
	result := make(map[string]interface{}, len(questionNames))

	questionMap := q[quiz]
	if questionMap == nil {
		// not even question map is present
		// return blank map
		return result
	}

	for _, questionName := range questionNames {
		question := questionMap[questionName]
		if question == nil {
			continue
		}
		result[questionName] = question.ToCounts()
	}
	return result
}
