# nre
Regex expression with natural language, powered by GPT.

## Install

```
go install github.com/yah01/nre@latest"
```

## Example
```
nre -d "errors.New(), and select the parameter string"
Regex: `errors\.New\((\".*\")\)` 

Explanation:
- `errors\.New\(` matches the literal string "errors.New("
- `(\".*\")` matches a string enclosed in double quotes, including any characters inside the quotes. The backslashes before the quotes are used to escape them, since quotes have a special meaning in regex. The parentheses around this part of the regex capture the string inside the quotes as a group, which can be accessed separately.
```