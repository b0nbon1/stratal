# Stratal Examples

This directory contains example job configurations that demonstrate various features of the Stratal task runner.

## Available Examples

### simple_parallel_example.json
Basic example of parallel execution with output passing between tasks.

### placeholder_to_email_example.json
**NEW**: Demonstrates fetching placeholder text from lorem ipsum API and sending it via email.

### placeholder_variations_example.json  
**NEW**: Shows different placeholder text options (short, medium, long) combined into a single email.

## New Builtin Tasks

### placeholder_text
Fetches lorem ipsum placeholder text from loripsum.net API.

**Parameters:**
- `paragraphs` (optional): Number of paragraphs to generate (default: 3)
- `sentences` (optional): Number of sentences per paragraph (default: 5)  
- `words` (optional): Generate specific number of words instead of paragraphs/sentences

**Examples:**
```json
{
  "name": "fetch_text",
  "type": "builtin", 
  "config": {
    "parameters": {
      "paragraphs": "2",
      "sentences": "4"
    }
  }
}
```

```json
{
  "name": "fetch_words",
  "type": "builtin",
  "config": {
    "parameters": {
      "words": "100"
    }
  }
}
```

### send_email (Updated)
Sends emails via SMTP with improved context handling and output formatting.

**Required Parameters:**
- `smtp_host`: SMTP server hostname
- `smtp_port`: SMTP server port  
- `smtp_user`: SMTP username
- `smtp_password`: SMTP password or app password
- `from`: Sender email address
- `to`: Recipient email address(es), comma-separated
- `subject`: Email subject line

**Body Parameters (at least one required):**
- `body_html`: HTML email body
- `body_text`: Plain text email body

**Example:**
```json
{
  "name": "send_notification",
  "type": "builtin",
  "config": {
    "depends_on": ["fetch_placeholder_text"],
    "parameters": {
      "smtp_host": "smtp.gmail.com",
      "smtp_port": "587", 
      "smtp_user": "user@gmail.com",
      "smtp_password": "app-password",
      "from": "user@gmail.com",
      "to": "recipient@example.com",
      "subject": "Generated Content",
      "body_text": "${TASK_OUTPUT_FETCH_PLACEHOLDER_TEXT}",
      "body_html": "<h1>Content</h1><p>${TASK_OUTPUT_FETCH_PLACEHOLDER_TEXT}</p>"
    }
  }
}
```

## Task Output Variables

Task outputs can be referenced in subsequent tasks using the format:
`${TASK_OUTPUT_TASK_NAME_IN_UPPERCASE}`

For example, if a task is named `fetch_placeholder_text`, its output can be accessed as:
`${TASK_OUTPUT_FETCH_PLACEHOLDER_TEXT}`

## Usage Tips

1. **Gmail SMTP**: Use app passwords instead of your regular password for Gmail SMTP
2. **HTML Formatting**: Use HTML in `body_html` for rich email formatting
3. **Dependencies**: Use `depends_on` to ensure tasks run in the correct order
4. **Parallel Execution**: Tasks with the same `order` value run in parallel
5. **Error Handling**: Both email and placeholder tasks include proper error handling and context cancellation support

## Running Examples

To test these examples, update the email credentials in the JSON files and use the Stratal API or CLI to submit the job configurations. 