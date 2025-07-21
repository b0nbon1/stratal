package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// FormatOutputTask formats and prints data with user-specified formatting
func FormatOutputTask(ctx context.Context, params map[string]string) (string, error) {
	// Get the template - this is the main formatting template
	template, exists := params["template"]
	if !exists || template == "" {
		return "", fmt.Errorf("missing required parameter: template")
	}

	// Get format type (default: text)
	formatType := strings.ToLower(params["format"])
	if formatType == "" {
		formatType = "text"
	}

	// Get any additional data to include
	data := make(map[string]interface{})

	// Add current timestamp if requested
	if params["include_timestamp"] == "true" {
		timestampFormat := params["timestamp_format"]
		if timestampFormat == "" {
			timestampFormat = "2006-01-02 15:04:05"
		}
		data["timestamp"] = time.Now().Format(timestampFormat)
		data["iso_timestamp"] = time.Now().Format(time.RFC3339)
	}

	// Add any custom fields
	for key, value := range params {
		if strings.HasPrefix(key, "field_") {
			fieldName := strings.TrimPrefix(key, "field_")
			data[fieldName] = value
		}
	}

	// Process the template based on format type
	var result string
	var err error

	switch formatType {
	case "json":
		result, err = formatAsJSON(template, data, params)
	case "table":
		result, err = formatAsTable(template, data, params)
	case "csv":
		result, err = formatAsCSV(template, data, params)
	case "xml":
		result, err = formatAsXML(template, data, params)
	case "text", "":
		result, err = formatAsText(template, data, params)
	default:
		return "", fmt.Errorf("unsupported format type: %s", formatType)
	}

	if err != nil {
		return "", fmt.Errorf("formatting error: %w", err)
	}

	return result, nil
}

// formatAsText processes a text template with variable substitution
func formatAsText(template string, data map[string]interface{}, params map[string]string) (string, error) {
	result := template

	// Replace data variables
	for key, value := range data {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	// Replace task output variables (these come from parameter interpolation)
	// Look for ${TASK_OUTPUT.task_name} or ${task_name.output} patterns
	taskOutputPattern := regexp.MustCompile(`\$\{([^}]+)\.output\}|\$\{TASK_OUTPUT\.([^}]+)\}`)
	matches := taskOutputPattern.FindAllStringSubmatch(result, -1)

	for _, match := range matches {
		fullMatch := match[0]
		var taskName string
		if match[1] != "" {
			taskName = match[1]
		} else {
			taskName = match[2]
		}

		// Look for the task output in params (it should be available as TASK_OUTPUT_TASKNAME)
		taskOutputKey := fmt.Sprintf("TASK_OUTPUT_%s", strings.ToUpper(strings.ReplaceAll(taskName, "-", "_")))
		if output, exists := params[taskOutputKey]; exists {
			result = strings.ReplaceAll(result, fullMatch, output)
		}
	}

	// Handle special formatting
	if params["uppercase"] == "true" {
		result = strings.ToUpper(result)
	}
	if params["lowercase"] == "true" {
		result = strings.ToLower(result)
	}

	// Handle line breaks
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\t", "\t")

	return result, nil
}

// formatAsJSON creates a JSON output with the template as structure
func formatAsJSON(template string, data map[string]interface{}, params map[string]string) (string, error) {
	// Try to parse template as JSON first
	var jsonData interface{}
	if err := json.Unmarshal([]byte(template), &jsonData); err != nil {
		// If not valid JSON, create a simple structure
		jsonData = map[string]interface{}{
			"message": template,
			"data":    data,
		}
	}

	// Add metadata if requested
	if params["include_metadata"] == "true" {
		metadata := map[string]interface{}{
			"generated_at": time.Now().Format(time.RFC3339),
			"format_type":  "json",
		}

		switch v := jsonData.(type) {
		case map[string]interface{}:
			v["_metadata"] = metadata
		default:
			jsonData = map[string]interface{}{
				"content":   jsonData,
				"_metadata": metadata,
			}
		}
	}

	// Pretty print if requested
	if params["pretty"] == "true" {
		bytes, err := json.MarshalIndent(jsonData, "", "  ")
		return string(bytes), err
	}

	bytes, err := json.Marshal(jsonData)
	return string(bytes), err
}

// formatAsTable creates a simple table format
func formatAsTable(template string, data map[string]interface{}, params map[string]string) (string, error) {
	var result strings.Builder

	// Table header
	separator := params["separator"]
	if separator == "" {
		separator = " | "
	}

	// Use template as header if it contains separators
	if strings.Contains(template, separator) {
		result.WriteString(template + "\n")
		// Add separator line
		headerLength := len(template)
		result.WriteString(strings.Repeat("-", headerLength) + "\n")
	} else {
		result.WriteString("Field" + separator + "Value\n")
		result.WriteString("-----" + separator + "-----\n")
	}

	// Add data rows
	for key, value := range data {
		result.WriteString(fmt.Sprintf("%s%s%v\n", key, separator, value))
	}

	return result.String(), nil
}

// formatAsCSV creates CSV output
func formatAsCSV(template string, data map[string]interface{}, params map[string]string) (string, error) {
	var result strings.Builder

	delimiter := params["delimiter"]
	if delimiter == "" {
		delimiter = ","
	}

	// CSV header from template or default
	headers := strings.Split(template, delimiter)
	if len(headers) <= 1 {
		headers = []string{"field", "value"}
	}

	result.WriteString(strings.Join(headers, delimiter) + "\n")

	// Add data
	for key, value := range data {
		row := []string{key, fmt.Sprintf("%v", value)}
		result.WriteString(strings.Join(row, delimiter) + "\n")
	}

	return result.String(), nil
}

// formatAsXML creates simple XML output
func formatAsXML(template string, data map[string]interface{}, params map[string]string) (string, error) {
	var result strings.Builder

	rootElement := params["root_element"]
	if rootElement == "" {
		rootElement = "data"
	}

	result.WriteString(fmt.Sprintf("<%s>\n", rootElement))

	if template != "" {
		result.WriteString(fmt.Sprintf("  <message>%s</message>\n", template))
	}

	for key, value := range data {
		result.WriteString(fmt.Sprintf("  <%s>%v</%s>\n", key, value, key))
	}

	result.WriteString(fmt.Sprintf("</%s>\n", rootElement))

	return result.String(), nil
}
