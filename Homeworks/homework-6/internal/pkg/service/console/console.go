package console

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"unicode"
)

var (
	ErrEmptyWord = errors.New("word is empty")
)

func New() *Service {
	help := map[string]*map[string]string{
		"Console": {
			"Help":  "Получить информацию обо всех имеющихся командах ",
			"Spell": "Вывести в консоль все буквы переданного на вход слова через пробел",
			"GoFmt": "На вход принимает *.txt файл, на выходе перед каждым абзацем вставляет таб и ставит точку в конце предложений",
		},
		"Solution": {
			"Get": "Получить статус по решению", "Create": "Послать решение студента",
		},
		"Student": {
			"Create": "Создать профиль студента",
			"Get":    "Получить информацию о студенте по его id",
			"Update": "Обновить информацию о студенте",
		},
		"Task": {
			"Get": "Получить список текущих задач",
		},
	}
	return &Service{
		help: help,
	}
}

type Service struct {
	help map[string]*map[string]string
}

// Help для получения информации обо всех имеющихся консольных командах
func (s *Service) Help() {
	for serviceName, methods := range s.help {
		fmt.Println(serviceName)
		for method, description := range *methods {
			fmt.Println(method, ":", description)
		}
		fmt.Println()
	}
}

func (s *Service) Spell(word string) error {
	if word == "" {
		return ErrEmptyWord
	}
	for _, rune := range word {
		fmt.Print(string(rune), " ")
	}
	fmt.Println()
	return nil
}

func (s *Service) GoFmt(file multipart.File) ([]byte, error) {
	//заведём билдер для оптимизированного добавления строк
	var builder strings.Builder
	_, err := io.Copy(&builder, file)
	if err != nil {
		return nil, err
	}

	contents := builder.String()
	builder.Reset()
	paragraphs := strings.Split(contents, "\r\n")
	for _, paragraph := range paragraphs {
		// табуляция к новому абзацу
		builder.WriteString("\t")
		words := strings.Split(paragraph, " ")
		isFirstWord := true
		for _, word := range words {
			if isFirstWord {
				isFirstWord = false
				builder.WriteString(word)
			} else {
				// точка перед словами, начинающихся с заглавной буквы
				if unicode.IsUpper([]rune(word)[0]) {
					builder.WriteString(".")
				}
				builder.WriteString(" ")
				builder.WriteString(word)
			}
		}
		builder.WriteString(".")
		builder.WriteString("\n")
	}

	return []byte(builder.String()), nil
}
