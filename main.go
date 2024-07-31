package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	TAG_SECURITY = "SECURITY"
	TAG_DOCUMENT = "DOCUMENT"
	TAG_TEST     = "TEST"
)

func main() {

	fmt.Println("SECRET:", os.Getenv("OPENAI_API_KEY"))
	fmt.Println("variable:", os.Getenv("APIKEY"))
	CatchApiKeyOpenAI()
	// Get the hash of the current commit
	commitHash := os.Getenv("GITHUB_SHA")
	fmt.Println("Hash do commit:", commitHash)

	// Get changes from commit
	cmd := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", commitHash)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting changes from commit:", err)
		os.Exit(1)
	}

	// Separate changed file names in a slice
	changedFiles := strings.Split(string(output), "\n")

	// Search and request units for ChatGPT, .pas only
	for _, file := range changedFiles {
		fmt.Printf("Testing extension : %v\n", file)
		if strings.HasSuffix(file, ".pas") {
			fmt.Printf("Searching Tag : %v\n", file)
			for ExistTags(file) {
				err := ProcessInDelphiFile(file)
				if err != nil {
					fmt.Printf("Error processing the file %s: %v\n", file, err)
				} else {
					fmt.Printf("File processed: %s\n", file)
				}
			}
		}
	}
}

func ExistTags(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err.Error())
		return false
	}

	// Check for the comment pattern in Delphi (//)
	pattern := "(//<" + TAG_SECURITY + ">|//<" + TAG_TEST + ">|//<" + TAG_DOCUMENT + ">)"
	match, err := regexp.Match(pattern, content)

	if err != nil {
		fmt.Printf("Error when using regex : %v\n", err.Error())
		return false
	}

	fmt.Printf("Result Regex : %v\n", match)
	return match

}

func ExtractCodeTag(content string, tag string) (string, error) {
	var result string

	// Sets the opening and closing comment pattern
	startComment := fmt.Sprintf("//<%s>", tag)
	endComment := fmt.Sprintf("//</%s>", tag)

	fmt.Printf("Extraindo tag : %v\n", startComment)

	// Assemble the regular expression with capturing groups
	pattern := fmt.Sprintf(`(?s)%s(.*?)%s`, regexp.QuoteMeta(startComment), regexp.QuoteMeta(endComment))
	regex := regexp.MustCompile(pattern)

	// Finds the first match in the text
	match := regex.FindStringSubmatch(content)

	// Checks if the tag was found
	if len(match) >= 2 {
		// Extracts captured text (code block)
		code := match[1]
		return code, nil
	}
	return result, nil
}

func FetchCodeFirstTag(code string) (string, string, string) {
	var result string
	var tag string
	var action string

	result, _ = ExtractCodeTag(code, TAG_DOCUMENT)
	tag = TAG_DOCUMENT
	action = "Create the following source comment and return the comment to me without accents in the comment writing: "
	if result == "" {
		result, _ = ExtractCodeTag(code, TAG_TEST)
		tag = TAG_TEST
		action = "Create a unit test method for the following source: "
	}
	if result == "" {
		result, _ = ExtractCodeTag(code, TAG_SECURITY)
		tag = TAG_SECURITY
		action = "Carry out a security analysis on the source below and send me a comment back with the security improvements without accents in the comment writing: "
	}

	if result == "" {
		panic("No codes found")
	}

	return result, tag, action
}

func CatchApiKeyOpenAI() string {
	result := os.Getenv("OPENAI_API_KEY")
	if result == "" {
		println("ApiKey:" + result)
		os.Exit(1)
		panic("Error: API key not found. Check if the secret OPENAI_API_KEY is configured.")
	}
	return result
}

func GetIndentation(code string, tagUsada string) string {
	startTag := fmt.Sprintf("//<%s>", tagUsada)
	endTag := fmt.Sprintf("//</%s>", tagUsada)
	startIndex := strings.Index(string(code), startTag)
	endIndex := strings.Index(string(code), endTag)

	// Check if tags were found
	if startIndex != -1 && endIndex != -1 {
		// Get the indentation of the original code block
		indentation := ""
		lines := strings.Split(code, "\n")
		if len(lines) > 0 {
			// Find the space or tab at the beginning of the first line
			for _, ch := range lines[0] {
				if ch == ' ' || ch == '\t' {
					indentation += string(ch)
				} else {
					break
				}
			}
			return indentation
		}
	}
	return ""
}

func ProcessInDelphiFile(filename string) error {
	var TagUsada string
	var code string
	var action string

	// Read file contents
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file : %v\n", err.Error())
	}

	// Search first tag code block found
	code, TagUsada, action = FetchCodeFirstTag(string(content))

	//Get token from chatGPT Api
	apiKey := CatchApiKeyOpenAI()

	prompt := action + string(code) // Use the file contents as a prompt with what action chatGPT should take

	fmt.Printf("Processing tag : %v\n", TagUsada)

	//Montando Body da requisição
	data := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  2048,               // Set the desired maximum number of tokens
		"model":       "text-davinci-003", // Specify the desired model here
		"temperature": 0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//request
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	// Call in API da OpenAI
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//result api
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Resposta da api : %v\n", string(responseData))

	// error return
	var responseError struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(responseData, &responseError)
	if err != nil {
		return err
	}

	// If the API returns any errors, deal here
	if responseError.Error.Message != "" {
		return errors.New(responseError.Error.Message)
	}

	// Extract the "text" tag from the response
	var responseText struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}

	err = json.Unmarshal(responseData, &responseText)
	if err != nil {
		return err
	}
	// Check chatGPT response for choices
	if len(responseText.Choices) > 0 {
		// Get the text of the first choice
		text := responseText.Choices[0].Text
		fmt.Println("Text return:", text)

		// Finding the position of tags in the original content
		startTag := fmt.Sprintf("//<%s>", TagUsada)
		endTag := fmt.Sprintf("//</%s>", TagUsada)
		startIndex := strings.Index(string(content), startTag)
		endIndex := strings.Index(string(content), endTag)
		indentation := GetIndentation(code, TagUsada)

		// Apply indentation and comments to the text returned by the API
		indentedText := ""
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if line != "" {
				if TagUsada == TAG_TEST {
					indentedText += indentation + "//" + line + "\n"
				} else {
					indentedText += indentation + line + "\n"
				}
			}
		}

		//Adding code sent to the api in the slice so that it is not removed from the file.
		indentedText += code

		// Check if tags were found
		if startIndex != -1 && endIndex != -1 {
			// Create a new byte slice for the new content
			newContent := make([]byte, 0, len(content)+len(indentedText)-len(code))

			// Copy content before first tag
			newContent = append(newContent, content[:startIndex]...)

			// Copy the text returned by the API
			newContent = append(newContent, []byte(indentedText)...)

			// Copy content after last tag
			newContent = append(newContent, content[endIndex+len(endTag):]...)

			fmt.Println("Engraved text:", string(newContent))

			// Write modified content back to file
			err = ioutil.WriteFile(filename, newContent, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return errors.New("tags not found in file")
		}
	} else {
		return errors.New("No choice found in answer")
	}

	return nil
}
