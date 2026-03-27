package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func IsAbsoluteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func IsRelativePath(path string) bool {
	if IsAbsoluteURL(path) {
		return false
	}
	if strings.HasPrefix(path, "/") {
		return false
	}
	if strings.Contains(path, "@") && !strings.Contains(path, "/") {
		return false
	}
	if strings.HasPrefix(path, "#") {
		return false
	}
	return true
}

func ConvertRelativeLinks(markdownContent, githubURL, branch string) string {
	githubURL = strings.TrimRight(githubURL, "/")

	// Pattern 1: Complex image links with clickable image [![alt](path)](path)
	reComplex := regexp.MustCompile(`\[!\[([^\]]*)\]\(([^)]+)\)\]\(([^)]+)\)`)
	markdownContent = reComplex.ReplaceAllStringFunc(markdownContent, func(match string) string {
		submatches := reComplex.FindStringSubmatch(match)
		if len(submatches) != 4 {
			return match
		}
		altText := submatches[1]
		imagePath := submatches[2]
		linkPath := submatches[3]

		if IsRelativePath(imagePath) {
			imagePath = fmt.Sprintf("%s/%s", githubURL, imagePath)
		}
		if IsRelativePath(linkPath) {
			linkPath = fmt.Sprintf("%s/%s", githubURL, linkPath)
		}
		return fmt.Sprintf("[![%s](%s)](%s)", altText, imagePath, linkPath)
	})

	// Pattern 2: Simple image links ![alt](path)
	reSimpleImage := regexp.MustCompile(`(?<!\[)!\[([^\]]*)\]\(([^)]+)\)(?!\])`)
	markdownContent = reSimpleImage.ReplaceAllStringFunc(markdownContent, func(match string) string {
		submatches := reSimpleImage.FindStringSubmatch(match)
		if len(submatches) != 3 {
			return match
		}
		altText := submatches[1]
		imagePath := submatches[2]

		if IsRelativePath(imagePath) {
			imagePath = fmt.Sprintf("%s/%s", githubURL, imagePath)
			return fmt.Sprintf("![%s](%s)", altText, imagePath)
		}
		return match
	})

	// Pattern 3: Regular links [text](path)
	reRegularLink := regexp.MustCompile(`(?<!\!)(?<!\])\[([^\]]+)\]\(([^)]+)\)`)
	markdownContent = reRegularLink.ReplaceAllStringFunc(markdownContent, func(match string) string {
		submatches := reRegularLink.FindStringSubmatch(match)
		if len(submatches) != 3 {
			return match
		}
		linkText := submatches[1]
		linkPath := submatches[2]

		if IsRelativePath(linkPath) {
			linkPath = fmt.Sprintf("%s/%s", githubURL, linkPath)
			return fmt.Sprintf("[%s](%s)", linkText, linkPath)
		}
		return match
	})

	return markdownContent
}

func FetchGithubData(githubURL string) (map[string]interface{}, map[string]interface{}, error) {
	re := regexp.MustCompile(`https?://github\.com/([^/]+)/([^/]+)(?:/.*)?$`)
	matches := re.FindStringSubmatch(githubURL)
	
	if len(matches) != 3 {
		return nil, nil, fmt.Errorf("invalid GitHub URL format")
	}
	
	username := matches[1]
	repoName := matches[2]
	
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repoName)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, nil, err
	}
	
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("github API returned status: %d", resp.StatusCode)
	}
	
	body, _ := io.ReadAll(resp.Body)
	var githubData map[string]interface{}
	if err := json.Unmarshal(body, &githubData); err != nil {
		return nil, nil, err
	}
	
	// Fetch README.md
	rawGithubURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/", username, repoName)
	readmeURL := rawGithubURL + "README.md"
	
	readmeResp, err := client.Get(readmeURL)
	if err == nil && readmeResp.StatusCode == http.StatusOK {
		readmeBody, _ := io.ReadAll(readmeResp.Body)
		githubData["readme_file"] = ConvertRelativeLinks(string(readmeBody), rawGithubURL, "main")
		readmeResp.Body.Close()
	} else {
		githubData["readme_file"] = nil
		if err == nil {
			readmeResp.Body.Close()
		}
	}
	
	// Fetch languages
	if languagesURL, ok := githubData["languages_url"].(string); ok && languagesURL != "" {
		langReq, _ := http.NewRequest("GET", languagesURL, nil)
		langReq.Header = req.Header // Copy auth headers
		langResp, err := client.Do(langReq)
		if err == nil && langResp.StatusCode == http.StatusOK {
			langBody, _ := io.ReadAll(langResp.Body)
			var langData map[string]interface{}
			json.Unmarshal(langBody, &langData)
			githubData["languages"] = langData
			langResp.Body.Close()
		} else {
			githubData["languages"] = map[string]interface{}{}
			if err == nil {
				langResp.Body.Close()
			}
		}
	} else {
		githubData["languages"] = map[string]interface{}{}
	}
	
	name, _ := githubData["name"].(string)
	description, _ := githubData["description"].(string)
	
	basicData := map[string]interface{}{
		"title":       name,
		"description": description,
		"url":         githubURL,
	}
	
	return basicData, githubData, nil
}
