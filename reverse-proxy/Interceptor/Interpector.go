package Interceptor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type InterceptionCondition struct {
	Includes []string
	//	 this may take some time develop and well test

	// NotIncludes []string
	// Regex       []string
}

type InterceptRequest struct {
	Method  []string
	Headers InterceptionCondition
	Query   InterceptionCondition
	Path    InterceptionCondition
	Body    InterceptionCondition
}

func isMethodFlagged(r *http.Request, condition InterceptRequest) bool {
	if len(condition.Method) == 0 {
		return true
	} else {
		for _, method := range condition.Method {
			if strings.ToUpper(r.Method) == strings.ToUpper(method) {
				return true
			}
		}
	}
	return false
}
func isHeaderFlagged(r *http.Request, condition InterceptRequest) bool {

	if len(condition.Headers.Includes) == 0 {
		return true
	} else {

		for _, header := range condition.Headers.Includes {
			if r.Header.Get(header) != "" {
				return true
			}
		}
	}
	// for _, header := range condition.Headers.NotIncludes {
	// 	if r.Header.Get(header) == "" {
	// 		return true
	// 	}
	// }
	// for _, pattern := range condition.Headers.Regex {
	// 	matched, err := regexp.MatchString(pattern, r.Header.Get(pattern))
	// 	if err != nil || !matched {
	// 		return true
	// 	}
	// }
	return false
}
func isQueryParamsFlagged(r *http.Request, condition InterceptRequest) bool {
	queryValues := r.URL.Query()
	if len(condition.Query.Includes) == 0 {
		return true
	} else {
		for _, param := range condition.Query.Includes {

			if _, ok := queryValues[param]; ok {

				return true
			}
		}
	}
	// for _, param := range condition.Query.NotIncludes {
	// 	if _, ok := queryValues[param]; ok {
	// 		return true, nil
	// 	}
	// }

	// for _, pattern := range condition.Query.Regex {
	// 	matched, err := regexp.MatchString(pattern, r.URL.RawQuery)
	// 	if err != nil || !matched {
	// 		return true, err
	// 	}
	// }
	return false
}
func isPathFlagged(r *http.Request, condition InterceptRequest) bool {

	if len(condition.Path.Includes) == 0 {
		return true
	} else {
		for _, path := range condition.Path.Includes {
			if strings.Contains(r.URL.Path, path) {
				return true
			}
		}
	}
	// for _, path := range condition.Path.NotIncludes {
	// 	if strings.Contains(r.URL.Path, path) {
	// 		return true, nil
	// 	}
	// }
	// for _, pattern := range condition.Path.Regex {
	// 	matched, err := regexp.MatchString(pattern, r.URL.Path)
	// 	if err != nil || !matched {
	// 		return true, err
	// 	}
	// }
	return false
}
func isBodyFlagged(r *http.Request, condition InterceptRequest) bool {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return true
	}

	defer r.Body.Close()
	if len(condition.Body.Includes) == 0 {
		return true
	}
	// Unmarshal the body into a map
	var bodyContent map[string]interface{}
	err = json.Unmarshal(body, &bodyContent)
	if err != nil {
		return true
	}

	// Validate the request body
	for _, value := range condition.Body.Includes {
		//Check if the value exists in the bodyContent map
		if _, exists := bodyContent[value]; exists {
			return true
		}
	}
	// for _, value := range condition.Body.NotIncludes {
	// 	//Check if the value does not exist in the bodyContent map
	// 	if _, exists := bodyContent[value]; exists {
	// 		return true, nil
	// 	}
	// }
	// for _, pattern := range condition.Body.Regex {
	// 	//Check if the body matches the regex pattern

	// 	re, err := regexp.Compile(pattern)
	// 	if err != nil {
	// 		return true, err
	// 	}
	// 	matched := false
	// 	for key := range bodyContent {
	// 		if re.MatchString(key) {
	// 			matched = true
	// 			break
	// 		}
	// 	}
	// 	if !matched {
	// 		return true, nil
	// 	}
	// }
	return false
}
func isRequestFlaggedOnSingleCondition(r *http.Request, condition InterceptRequest) bool {

	MethodFlagged := false
	HeaderFlagged := false
	QueryParamsFlagger := false
	PathFlagged := false
	bodyFlagged := false
	// Validate the request method

	if flagged := isMethodFlagged(r, condition); flagged {
		MethodFlagged = true
	}
	// Validate the request headers
	if flagged := isHeaderFlagged(r, condition); flagged {
		HeaderFlagged = true
	}
	// Validate the request query parameters

	if flagged := isQueryParamsFlagged(r, condition); flagged {
		QueryParamsFlagger = true
	}
	// Validate the request path

	if flagged := isPathFlagged(r, condition); flagged {
		PathFlagged = true
	}

	if flagged := isBodyFlagged(r, condition); flagged {
		bodyFlagged = true
	}
	// If none of the validation checks failed, the request shall pass
	return MethodFlagged && HeaderFlagged && QueryParamsFlagger && PathFlagged && bodyFlagged
}
func IsRequestFlagged(r *http.Request, conditions []InterceptRequest) bool {
	RequestParamsAreFlagged := false
	for _, condition := range conditions {
		value := isRequestFlaggedOnSingleCondition(r, condition)
		if value {
			RequestParamsAreFlagged = value
			break
		}
	}
	return RequestParamsAreFlagged
}
